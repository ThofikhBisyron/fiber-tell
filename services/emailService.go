package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type EmailService struct {
	apiKey      string
	senderEmail string
	senderName  string
}

func NewEmailService() *EmailService {
	return &EmailService{
		apiKey:      os.Getenv("BREVO_API"),
		senderEmail: os.Getenv("BREVO_SENDER_EMAIL"),
		senderName:  os.Getenv("BREVO_SENDER_NAME"),
	}

}

func (s *EmailService) SendOtp(
	email string,
	code string,
) error {
	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  s.senderName,
			"email": s.senderEmail,
		},
		"to": []map[string]string{
			{
				"email": email,
			},
		},
		"subject": "Go Tell OTP Code",
		"htmlContent": fmt.Sprintf(
			`
			<h2> GG TELL <h2>
			<p> Your OTP Code is: <p>
			<h1>%s</h1>
			<p>This Code will expire in 5 minutes<p>
			`, code),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", s.apiKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"brevo returned status code %d",
			resp.StatusCode,
		)
	}

	return nil
}
