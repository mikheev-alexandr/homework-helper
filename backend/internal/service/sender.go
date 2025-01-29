package service

import (
	"os"

	"gopkg.in/gomail.v2"
)

type SMTPEmailSender struct {
	Host string
	Port int
}

func NewEmailSenderService(host string, port int) *SMTPEmailSender {
	return &SMTPEmailSender{
		Host: host,
		Port: port,
	}
}

func (s *SMTPEmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()

	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	m.SetHeader("From", "youremail@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, email, password)

	return d.DialAndSend(m)
}
