package otp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type Service interface {
	SendOTP(db *gorm.DB, phoneNumber string, purpose string, dropOffID uint) (*OTPRequest, error)
	VerifyOTP(db *gorm.DB, phoneNumber string, code, purpose string) (bool, error)
	SendEmailOTP(emailAddress, code string) (*EmailResponse, error)
	SendWhatsAppOTP(phoneNumber string, data map[string]string) (*WhatsAppResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type EmailRequest struct {
	APIKey               string `json:"api_key"`
	EmailAddress         string `json:"email_address"`
	Code                 string `json:"code"`
	EmailConfigurationID string `json:"email_configuration_id"`
}

type EmailResponse struct {
	Code      string `json:"code"`
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
	Balance   string `json:"balance"`
	User      string `json:"user"`
}
type WhatsAppRequest struct {
	PhoneNumber string            `json:"phone_number"`
	DeviceID    string            `json:"device_id"`
	TemplateID  string            `json:"template_id"`
	APIKey      string            `json:"api_key"`
	Data        map[string]string `json:"data"`
}

// WhatsAppResponse represents the response from the WhatsApp OTP send template endpoint.
type WhatsAppResponse struct {
	Code      string  `json:"code"`
	MessageID string  `json:"message_id"`
	Message   string  `json:"message"`
	Balance   float64 `json:"balance"`
	User      string  `json:"user"`
}

func (s *service) SendOTP(db *gorm.DB, phoneNumber string, purpose string, dropOffID uint) (*OTPRequest, error) {

	pinID, err := sendTermiiOTP(phoneNumber)
	if err != nil {
		return nil, err
	}

	otpRequest := &OTPRequest{
		PhoneNumber: phoneNumber,
		PinID:       pinID,
		Purpose:     purpose,
		DropOffID:   dropOffID,
	}

	if err := s.repo.SaveOTP(db, otpRequest); err != nil {
		return nil, err
	}

	return otpRequest, nil
}

func (s *service) VerifyOTP(db *gorm.DB, phoneNumber string, code string, purpose string) (bool, error) {

	otpRequest, err := s.repo.FindPinIDByPhoneNumber(db, phoneNumber, purpose)
	if err != nil {
		return false, errors.New("phone number not found")
	}

	return verifyTermiiOTP(otpRequest.PinID, code)
}

func sendTermiiOTP(phone string) (string, error) {
	apiKey := os.Getenv("TERMII_API_KEY")
	termiiBaseUrl := os.Getenv("TERMII_BASE_URL")
	if apiKey == "" {
		return "", errors.New("TERMII_API_KEY is not set")
	}
	if termiiBaseUrl == "" {
		return "", errors.New("TERMII_BASE_URL is not set")
	}

	// Transform phone from local format to international format
	// If phone starts with "0", replace it with "234" (Nigeria country code)
	formattedPhone := phone
	if len(phone) > 0 && phone[0] == '0' {
		formattedPhone = "234" + phone[1:]
	}

	payload := map[string]interface{}{
		"api_key":          apiKey,
		"message_type":     "ALPHANUMERIC",
		"to":               formattedPhone,
		"from":             "Child Safe",
		"channel":          "generic",
		"pin_attempts":     3,
		"pin_time_to_live": 2,
		"pin_length":       5,
		"pin_placeholder":  "< 123456 >",
		"message_text":     "Your ChildSafe code is < 123456 > it expires in two minutes",
		"pin_type":         "NUMERIC",
	}

	data, _ := json.Marshal(payload)
	sendURL := termiiBaseUrl + "/api/sms/otp/send"
	res, err := http.Post(sendURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return "", errors.New("failed to send OTP via Termii")
	}

	var resp struct {
		PinID string `json:"pinId"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", err
	}

	return resp.PinID, nil
}

func verifyTermiiOTP(pinID, code string) (bool, error) {
	apiKey := os.Getenv("TERMII_API_KEY")
	termiiBaseUrl := os.Getenv("TERMII_BASE_URL")
	if apiKey == "" {
		return false, errors.New("TERMII_API_KEY is not set")
	}
	if termiiBaseUrl == "" {
		return false, errors.New("TERMII_BASE_URL is not set")
	}

	payload := map[string]interface{}{
		"api_key": apiKey,
		"pin_id":  pinID,
		"pin":     code,
	}

	data, _ := json.Marshal(payload)
	verifyURL := termiiBaseUrl + "/api/sms/otp/verify"
	res, err := http.Post(verifyURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return false, fmt.Errorf("termii API returned status code: %d", res.StatusCode)
	}

	var termiiResponse struct {
		Verified          bool   `json:"verified"`
		PinID             string `json:"pinId"`
		Msisdn            string `json:"msisdn"`
		AttemptsRemaining int    `json:"attemptsRemaining"`
	}

	if err := json.NewDecoder(res.Body).Decode(&termiiResponse); err != nil {
		return false, fmt.Errorf("failed to decode termii response: %w", err)
	}

	return termiiResponse.Verified, nil
}

func (s *service) SendEmailOTP(emailAddress, code string) (*EmailResponse, error) {
	apiKey := os.Getenv("TERMII_API_KEY")
	termiiBaseURL := os.Getenv("TERMII_BASE_URL")
	emailConfigurationID := os.Getenv("TERMII_EMAIL_CONFIGURATION_ID")

	if apiKey == "" {
		return nil, errors.New("TERMII_API_KEY is not set")
	}
	if termiiBaseURL == "" {
		return nil, errors.New("TERMII_BASE_URL is not set")
	}
	if emailConfigurationID == "" {
		return nil, errors.New("TERMII_EMAIL_CONFIGURATION_ID is not set")
	}

	payload := EmailRequest{
		APIKey:               apiKey,
		EmailAddress:         emailAddress,
		Code:                 code,
		EmailConfigurationID: emailConfigurationID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	sendURL := termiiBaseURL + "/api/email/otp/send"
	res, err := http.Post(sendURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to Termii email OTP endpoint: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		bodyBytes, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return nil, fmt.Errorf("termii email OTP API returned status %d, failed to read error body: %v", res.StatusCode, readErr)
		}
		return nil, fmt.Errorf("termii email OTP API returned status %d, body: %s", res.StatusCode, string(bodyBytes))
	}

	var response EmailResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode Termii email OTP response: %w", err)
	}

	return &response, nil
}

func (s *service) SendWhatsAppOTP(phoneNumber string, data map[string]string) (*WhatsAppResponse, error) {
	apiKey := os.Getenv("TERMII_API_KEY")
	termiiBaseURL := os.Getenv("TERMII_BASE_URL")
	whatsappDeviceID := os.Getenv("TERMII_WHATSAPP_DEVICE_ID") // Assuming you'll set a specific device ID for WhatsApp
	templateID := os.Getenv("TERMII_WHATSAPP_TEMPLATE_ID")

	if apiKey == "" {
		return nil, errors.New("TERMII_API_KEY is not set")
	}
	if termiiBaseURL == "" {
		return nil, errors.New("TERMII_BASE_URL is not set")
	}
	if whatsappDeviceID == "" {
		return nil, errors.New("TERMII_WHATSAPP_DEVICE_ID is not set")
	}
	if templateID == "" {
		return nil, errors.New("templateID cannot be empty")
	}
	if phoneNumber == "" {
		return nil, errors.New("phoneNumber cannot be empty")
	}
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	// Format phone number to international format if it starts with '0' (Nigeria specific)
	formattedPhone := phoneNumber
	if len(phoneNumber) > 0 && phoneNumber[0] == '0' {
		formattedPhone = "234" + phoneNumber[1:]
	}

	payload := WhatsAppRequest{
		PhoneNumber: formattedPhone,
		DeviceID:    whatsappDeviceID,
		TemplateID:  templateID,
		APIKey:      apiKey,
		Data:        data,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	sendURL := termiiBaseURL + "/api/send/template"
	res, err := http.Post(sendURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to Termii WhatsApp template endpoint: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		bodyBytes, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return nil, fmt.Errorf("termii whatsapp template API returned status %d, failed to read error body: %v", res.StatusCode, readErr)
		}
		return nil, fmt.Errorf("termii whatsapp template API returned status %d, body: %s", res.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var singleResponse WhatsAppResponse
	err = json.Unmarshal(bodyBytes, &singleResponse)
	if err == nil {
		return &singleResponse, nil
	}

	// var arrayResponse WhatsAppResponse
	// err = json.Unmarshal(bodyBytes, &arrayResponse)
	// if err == nil && len(arrayResponse) > 0 {
	// 	return &arrayResponse[0], nil
	// }

	return nil, fmt.Errorf("failed to decode Termii WhatsApp template response: %v, body: %s", err, string(bodyBytes))
}
