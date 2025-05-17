package messaging

import (
	"fmt"
	"strings"
)

type service struct {
	emailClient *SendPulseClient
}

type Service interface {
	SendDropSessionEmail(payload DropSessionEmailPayload) error
}

func NewService(emailClient *SendPulseClient) Service {
	return &service{
		emailClient: emailClient,
	}
}

func (s *service) SendDropSessionEmail(payload DropSessionEmailPayload) error {
	children := strings.Join(payload.Children, ", ")
	html := fmt.Sprintf(`
		<h3>Dear %s,</h3>
		<p>Your children (%s) have been successfully dropped off at %s on %s.</p>
		<p>Your pickup secret is: <strong>%s</strong></p>
		<p>Please keep this code safe. You'll need it to pick them up.</p>
	`, payload.GuardianName, children, payload.ChurchName, payload.Date, payload.Secret)

	return s.emailClient.SendEmail(payload.GuardianEmail, "Drop-Off Confirmation & Pickup Secret", html)
}
