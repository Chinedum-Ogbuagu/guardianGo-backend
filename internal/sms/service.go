package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type Service interface {
	SendSMS(db *gorm.DB, phoneNumber, message string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SendSMS(db *gorm.DB, phoneNumber, message string) error {
	formattedPhone := phoneNumber
	if len(phoneNumber) > 0 && phoneNumber[0] == '0' {
		formattedPhone = "234" + phoneNumber[1:]
	}

	apiKey := os.Getenv("TERMII_API_KEY")
	termiiBaseUrl := os.Getenv("TERMII_BASE_URL")
	if apiKey == "" || termiiBaseUrl == "" {
		return errors.New("TERMII_API_KEY or TERMII_BASE_URL is not set")
	}

	payload := map[string]interface{}{
		"to":       formattedPhone,
		"from":     "ChildSafe",
		"sms":      message,
		"type":     "plain",
		"channel":  "generic",
		"api_key":  apiKey,
	}

	data, _ := json.Marshal(payload)
	sendURL := termiiBaseUrl + "/api/sms/send"
	res, err := http.Post(sendURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return errors.New("failed to send SMS")
	}

	// Optionally log the SMS
	log := &SMSLog{
		PhoneNumber: formattedPhone,
		Message:     message,
	}
	return s.repo.SaveSMSLog(db, log)
}
