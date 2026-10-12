package auth_service_test

import (
	"testing"
	"time"
)

type sentEmail struct {
	to      string
	subject string
	body    string
}

type mockMailer struct {
	sent chan sentEmail
	err  error
}

func newMockMailer() *mockMailer {
	return &mockMailer{sent: make(chan sentEmail, 10)}
}

func (m *mockMailer) Send(to, subject, body string) error {
	m.sent <- sentEmail{to: to, subject: subject, body: body}

	return m.err
}

func waitForEmail(t *testing.T, m *mockMailer) sentEmail {
	t.Helper()

	select {
	case email := <-m.sent:
		return email
	case <-time.After(2 * time.Second):
		t.Fatal("письмо не было отправлено за 2 секунды")

		return sentEmail{}
	}
}

func assertNoEmail(t *testing.T, m *mockMailer) {
	t.Helper()

	select {
	case email := <-m.sent:
		t.Fatalf("письмо не должно было отправляться, но ушло на %q", email.to)
	case <-time.After(100 * time.Millisecond):
	}
}
