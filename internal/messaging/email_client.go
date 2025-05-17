package messaging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SendPulseClient struct {
	AccessToken string
}

func NewSendPulseClient(token string) *SendPulseClient {
	return &SendPulseClient{
		AccessToken: token,
	}
}

func (c *SendPulseClient) SendEmail(to, subject, html string) error {
	body := map[string]interface{}{
		"email": map[string]interface{}{
			"to": []map[string]string{
				{"email": to},
			},
			"subject": subject,
			"from": map[string]string{
				"name":  os.Getenv("SENDPULSE_FROM_NAME"),
				"email": os.Getenv("SENDPULSE_FROM_EMAIL"),
			},
			"html": html,
		},
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://api.sendpulse.com/smtp/emails", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	return nil
}
