package core_mailer

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
)

type Mailer interface {
	Send(to, subject, body string) error
}

type smtpMailer struct {
	config Config
}

func NewMailer(config Config) Mailer {
	return &smtpMailer{config: config}
}

func (m *smtpMailer) Send(to, subject, body string) error {
	addr := net.JoinHostPort(m.config.Host, strconv.Itoa(m.config.Port))
	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	message := buildMessage(m.config.From, to, subject, body)

	if m.config.Port == 465 {
		return m.sendImplicitTLS(addr, auth, to, message)
	}

	if err := smtp.SendMail(addr, auth, m.config.From, []string{to}, message); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}

	return nil
}

func (m *smtpMailer) sendImplicitTLS(addr string, auth smtp.Auth, to string, message []byte) error {
	tlsConn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: m.config.Host})
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	defer tlsConn.Close()

	client, err := smtp.NewClient(tlsConn, m.config.Host)
	if err != nil {
		return fmt.Errorf("init smtp client: %w", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	if err := client.Mail(m.config.From); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	defer w.Close()

	if _, err := w.Write(message); err != nil {
		return fmt.Errorf("write message body: %w", err)
	}

	return nil
}

func buildMessage(from, to, subject, body string) []byte {
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		from, to, subject, body,
	)

	return []byte(msg)
}
