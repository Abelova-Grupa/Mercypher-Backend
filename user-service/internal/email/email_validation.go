package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"os"
)

type EmailAuth struct {
	FromEmail  string
	AppPass    string
	SmtpServer string
	Port       string
}

type EmailPayload struct {
	Username string `json:"username"`
	ToEmail string	`json:email`
	AuthCode string	`json:auth_code`
}

func NewEmailAuth() *EmailAuth {
	return &EmailAuth{
		FromEmail: os.Getenv("EMAIL_CLIENT"),
		AppPass: os.Getenv("EMAIL_APP_PASS"),
		SmtpServer: os.Getenv("SMTP_HOST"),
		Port: os.Getenv("SSL_SMTP_PORT"),
	}
}

var (
	ErrInvalidParams = errors.New("parameters are invalid")
	ErrInvalidEnvVars = errors.New("invalid env variables")
)

func (e *EmailAuth) SendEmail(payload EmailPayload) error{
	if e.AppPass == "" || e.FromEmail == "" || e.SmtpServer == "" || e.Port == "" {
		return ErrInvalidEnvVars
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         e.SmtpServer,
	}

	auth := smtp.PlainAuth("", e.FromEmail, e.AppPass, e.SmtpServer)
	conn, err := tls.Dial("tcp", e.SmtpServer+":"+e.Port, tlsConfig)
	if err != nil {
		return err
	}

	headers := make(map[string] string)
	headers["From"] = e.FromEmail
	headers["To"] = payload.ToEmail
	headers["Subject"] = "Verify your mercypher account"

	messageBody := "Your verification code is " + payload.AuthCode + ". You have 15 minutes to active your account."
	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n",k,v)
	}
	message += messageBody
	client, err := smtp.NewClient(conn, e.SmtpServer)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(e.FromEmail); err != nil {
		return err
	}

	if err := client.Rcpt(payload.ToEmail); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	_, err = wc.Write([]byte(message))
	if err != nil {
		return nil
	}
	wc.Close()
	return nil

}