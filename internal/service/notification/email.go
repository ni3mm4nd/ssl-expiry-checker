package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/domain/sslcheck"
)

type emailNotification struct {
	service SMTPConfig
}

type password string

type SMTPConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Host      string   `yaml:"host"`
	Port      int      `yaml:"port"`
	Password  password `yaml:"password"`
	Sender    string   `yaml:"sender"`
	Receivers []string `yaml:"receivers"`
	Subject   string   `yaml:"subject"`
	Message   string   `yaml:"message"`
	TLS       bool     `yaml:"tls"`
}

func (password) MarshalYAML() (interface{}, error) {
	return "******", nil
}

func (p password) String() string {
	return string(p)
}

func (e emailNotification) Notify(checks []sslcheck.SSLCheck) {
	log.Println("Email notification in progress...")

	from := e.service.Sender
	subject := e.service.Subject
	password := e.service.Password
	to := e.service.Receivers
	smtpHost := e.service.Host
	smtpPort := e.service.Port

	message := []byte(e.service.Message)

	auth := smtp.PlainAuth("", from, password.String(), smtpHost)

	msg := ""
	msg += fmt.Sprintf("From: %s\r\n", from)

	if len(to) > 0 {
		msg += fmt.Sprintf("To: %s\r\n", to)
	}

	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += fmt.Sprintf("\r\n%s\r\n", message)

	var s strings.Builder

	for _, c := range checks {
		s.WriteString(fmt.Sprintf("\n\nurl: %s\ndays left: %v days\n", c.TargetURL, c.DaysLeft))
		if c.Error != "" {
			s.WriteString(fmt.Sprintf("error: %s\n", c.Error))
		}
	}

	msg += fmt.Sprintf("\r\n%s\r\n", s.String())

	// Sending email.
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, to, []byte(msg))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func NewEmailNotification(config SMTPConfig) INotification {
	return &emailNotification{
		service: config,
	}
}
