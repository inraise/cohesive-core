package core_mailer

import (
	"encoding/base64"
	"mime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildMessage(t *testing.T) {
	const body = "Здравствуйте! Подтвердите email, перейдя по ссылке: https://example.com/verify?token=abc"

	message := string(buildMessage("noreply@example.com", "user@example.com", "Подтверждение email", body))

	head, encodedBody, found := strings.Cut(message, "\r\n\r\n")
	require.True(t, found, "заголовки должны отделяться от тела пустой строкой")

	t.Run("тема с кириллицей кодируется по RFC 2047", func(t *testing.T) {
		var subjectLine string

		unfolded := strings.ReplaceAll(head, "\r\n ", " ")

		for _, line := range strings.Split(unfolded, "\r\n") {
			if strings.HasPrefix(line, "Subject: ") {
				subjectLine = strings.TrimPrefix(line, "Subject: ")
			}
		}

		require.NotEmpty(t, subjectLine)
		assert.True(t, strings.HasPrefix(subjectLine, "=?utf-8?q?"), "получили %q", subjectLine)

		decoded, err := new(mime.WordDecoder).DecodeHeader(subjectLine)
		require.NoError(t, err)
		assert.Equal(t, "Подтверждение email", decoded)
	})

	t.Run("тело декодируется обратно без потерь", func(t *testing.T) {
		raw := strings.ReplaceAll(encodedBody, "\r\n", "")

		decoded, err := base64.StdEncoding.DecodeString(raw)
		require.NoError(t, err)
		assert.Equal(t, body, string(decoded))
	})

	t.Run("строки тела не длиннее 76 символов", func(t *testing.T) {
		for _, line := range strings.Split(encodedBody, "\r\n") {
			assert.LessOrEqual(t, len(line), 76)
		}
	})

	t.Run("в заголовках нет сырых не-ASCII символов", func(t *testing.T) {
		for _, r := range head {
			assert.Less(t, r, rune(128))
		}
	})
}

func TestParseAddress(t *testing.T) {
	t.Run("обычный адрес", func(t *testing.T) {
		addr, err := parseAddress("user@example.com")

		require.NoError(t, err)
		assert.Equal(t, "user@example.com", addr.Address)
	})

	t.Run("адрес с именем отдаёт чистый Address для конверта", func(t *testing.T) {
		addr, err := parseAddress("Cohesive <noreply@example.com>")

		require.NoError(t, err)
		assert.Equal(t, "noreply@example.com", addr.Address)
	})

	t.Run("перенос строки — попытка инъекции заголовка", func(t *testing.T) {
		_, err := parseAddress("user@example.com\r\nBcc: victim@example.com")

		require.Error(t, err)
	})

	t.Run("не email", func(t *testing.T) {
		_, err := parseAddress("not-an-email")

		require.Error(t, err)
	})
}
