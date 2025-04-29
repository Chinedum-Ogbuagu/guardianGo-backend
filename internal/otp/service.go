package otp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type Service interface {
    SendOTP(db *gorm.DB, phoneNumber string, purpose string, dropOffID uint) (*OTPRequest, error)
    VerifyOTP(db *gorm.DB, phoneNumber string, code, purpose string) (bool, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
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
        "api_key":            apiKey,
        "message_type":       "NUMERIC",
        "to":                 formattedPhone,
        "from":               "Child Safe",
        "channel":            "generic",
        "pin_attempts":       1,
        "pin_time_to_live":   20,
        "pin_length":         5,
        "pin_placeholder":    "< 123456 >",
        "message_text":       "Your ChildSafe code is < 123456 >",
        "pin_type":           "NUMERIC",
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
        return false, nil
    }

    return true, nil
}