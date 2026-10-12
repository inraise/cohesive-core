package core_mailer

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
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
	recipient, err := parseAddress(to)
	if err != nil {
		return fmt.Errorf("parse recipient: %w", err)
	}

	sender, err := parseAddress(m.config.From)
	if err != nil {
		return fmt.Errorf("parse sender: %w", err)
	}

	addr := net.JoinHostPort(m.config.Host, strconv.Itoa(m.config.Port))
	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	message := buildMessage(sender.String(), recipient.String(), subject, body)

	if m.config.Port == 465 {
		return m.sendImplicitTLS(addr, auth, sender.Address, recipient.Address, message)
	}

	if err = smtp.SendMail(addr, auth, sender.Address, []string{recipient.Address}, message); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}

	return nil
}

func (m *smtpMailer) sendImplicitTLS(
	addr string,
	auth smtp.Auth,
	from, to string,
	message []byte,
) error {
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

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}

	if _, err = w.Write(message); err != nil {
		return fmt.Errorf("write message body: %w", err)
	}

	// Сервер принимает или отклоняет письмо именно на Close — ошибку терять нельзя.
	if err = w.Close(); err != nil {
		return fmt.Errorf("close message body: %w", err)
	}

	// Письмо уже принято, ошибка QUIT на результат не влияет.
	_ = client.Quit()

	return nil
}

// parseAddress разбирает адрес и отсекает всё, что могло бы попасть в заголовок
// как перенос строки (инъекция заголовков через пользовательский email).
func parseAddress(raw string) (*mail.Address, error) {
	if strings.ContainsAny(raw, "\r\n") {
		return nil, errors.New("address contains line breaks")
	}

	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	return addr, nil
}

func buildMessage(from, to, subject, body string) []byte {
	headers := []string{
		"From: " + from,
		"To: " + to,
		"Subject: " + mime.QEncoding.Encode("utf-8", subject),
		"Date: " + time.Now().Format(time.RFC1123Z),
		"MIME-Version: 1.0",
		`Content-Type: text/plain; charset="utf-8"`,
		"Content-Transfer-Encoding: base64",
	}

	return []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + encodeBase64Body(body))
}

func encodeBase64Body(body string) string {
	const lineLen = 76

	encoded := base64.StdEncoding.EncodeToString([]byte(body))

	var sb strings.Builder

	for len(encoded) > lineLen {
		sb.WriteString(encoded[:lineLen])
		sb.WriteString("\r\n")

		encoded = encoded[lineLen:]
	}

	sb.WriteString(encoded)

	return sb.String()
}
