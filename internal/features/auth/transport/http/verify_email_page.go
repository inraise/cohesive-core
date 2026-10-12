package auth_transport_http

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"

	"go.uber.org/zap"
)

const (
	// Настоящий токен — 43 символа; всё, что заметно длиннее, заведомо мусор.
	maxVerificationTokenLen = 128
	maxConfirmFormBodyBytes = 4 << 10

	// Страница не грузит ничего внешнего и не выполняет скриптов.
	// form-action 'self' — форма может уходить только на наш же хост.
	verifyPageCSP = "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'; " +
		"frame-ancestors 'none'; base-uri 'none'"
)

type verifyPageKind string

const (
	verifyPageConfirm verifyPageKind = "confirm"
	verifyPageSuccess verifyPageKind = "success"
	verifyPageInvalid verifyPageKind = "invalid"
	verifyPageError   verifyPageKind = "error"
)

type verifyPageData struct {
	Kind  verifyPageKind
	Token string
}

var verifyPageTemplate = template.Must(template.New("verify_email").Parse(`<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="robots" content="noindex">
<title>Cohesive — подтверждение email</title>
<style>
:root { color-scheme: light dark; --bg:#f5f6f8; --card:#fff; --text:#1a1d21; --muted:#5c6570; --accent:#2f6fed; --ok:#1f9d55; --bad:#d64545; }
@media (prefers-color-scheme: dark) { :root { --bg:#14161a; --card:#1e2126; --text:#e8eaed; --muted:#9aa3ad; --accent:#5b8def; --ok:#3ecf7a; --bad:#ef6b6b; } }
* { box-sizing: border-box; }
body { margin:0; min-height:100vh; display:flex; align-items:center; justify-content:center; padding:24px;
       background:var(--bg); color:var(--text); font:16px/1.5 -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
main { width:100%; max-width:420px; background:var(--card); border-radius:16px; padding:36px 28px; text-align:center;
       box-shadow:0 4px 24px rgba(0,0,0,.08); }
.icon { width:56px; height:56px; margin:0 auto 18px; border-radius:50%; display:flex; align-items:center; justify-content:center; font-size:28px; color:#fff; }
.icon.confirm { background:var(--accent); } .icon.success { background:var(--ok); } .icon.invalid, .icon.error { background:var(--bad); }
h1 { font-size:22px; margin:0 0 10px; }
p { margin:0 0 22px; color:var(--muted); }
button { width:100%; padding:14px 18px; border:0; border-radius:10px; background:var(--accent); color:#fff; font:inherit; font-weight:600; cursor:pointer; }
button:hover { filter:brightness(1.08); }
</style>
</head>
<body>
<main>
{{if eq .Kind "confirm"}}
  <div class="icon confirm">&#9993;</div>
  <h1>Подтвердите email</h1>
  <p>Нажмите кнопку, чтобы завершить регистрацию в Cohesive.</p>
  <form method="post" action="verify-email/confirm">
    <input type="hidden" name="token" value="{{.Token}}">
    <button type="submit">Подтвердить email</button>
  </form>
{{else if eq .Kind "success"}}
  <div class="icon success">&#10003;</div>
  <h1>Готово!</h1>
  <p>Email подтверждён. Теперь можно войти в приложение.</p>
{{else if eq .Kind "invalid"}}
  <div class="icon invalid">!</div>
  <h1>Ссылка недействительна</h1>
  <p>Ссылка устарела или уже была использована. Запросите новое письмо в приложении.</p>
{{else}}
  <div class="icon error">!</div>
  <h1>Что-то пошло не так</h1>
  <p>Не удалось подтвердить email. Попробуйте ещё раз чуть позже.</p>
{{end}}
</main>
</body>
</html>
`))

func renderVerifyPage(
	rw http.ResponseWriter,
	log *core_logger.Logger,
	status int,
	data verifyPageData,
) {
	var buf bytes.Buffer
	if err := verifyPageTemplate.Execute(&buf, data); err != nil {
		log.Error("failed to render verify email page", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)

		return
	}

	header := rw.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Cache-Control", "no-store")
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Security-Policy", verifyPageCSP)

	rw.WriteHeader(status)

	if _, err := rw.Write(buf.Bytes()); err != nil {
		log.Error("write verify email page", zap.Error(err))
	}
}

// VerifyEmailPage godoc
// @Summary Страница подтверждения email
// @Description Открывается по ссылке из письма. Ничего не подтверждает, а только показывает кнопку: почтовые сканеры сами открывают ссылки из писем и не должны «сжигать» одноразовый токен
// @Tags auth
// @Produce html
// @Param token query string true "Токен из письма"
// @Success 200 {string} string "HTML-страница с кнопкой подтверждения"
// @Failure 400 {string} string "HTML-страница: ссылка недействительна"
// @Router /auth/verify-email [get]
func (h *AuthHTTPHandler) VerifyEmailPage(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())

	token := r.URL.Query().Get("token")
	if token == "" || len(token) > maxVerificationTokenLen {
		renderVerifyPage(rw, log, http.StatusBadRequest, verifyPageData{Kind: verifyPageInvalid})

		return
	}

	renderVerifyPage(rw, log, http.StatusOK, verifyPageData{Kind: verifyPageConfirm, Token: token})
}

// ConfirmVerifyEmail godoc
// @Summary Подтверждение email (форма со страницы)
// @Description Вызывается кнопкой на странице подтверждения. Отвечает HTML-страницей с результатом
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce html
// @Param token formData string true "Токен из письма"
// @Success 200 {string} string "HTML-страница: email подтверждён"
// @Failure 400 {string} string "HTML-страница: ссылка недействительна"
// @Failure 429 {object} cohesive-core_internal_core_transport_http_response.ErrorResponse "Too many requests"
// @Failure 500 {string} string "HTML-страница: внутренняя ошибка"
// @Router /auth/verify-email/confirm [post]
func (h *AuthHTTPHandler) ConfirmVerifyEmail(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	r.Body = http.MaxBytesReader(rw, r.Body, maxConfirmFormBodyBytes)

	if err := r.ParseForm(); err != nil {
		log.Warn("failed to parse verify email form", zap.Error(err))
		renderVerifyPage(rw, log, http.StatusBadRequest, verifyPageData{Kind: verifyPageInvalid})

		return
	}

	token := r.PostFormValue("token")
	if token == "" || len(token) > maxVerificationTokenLen {
		renderVerifyPage(rw, log, http.StatusBadRequest, verifyPageData{Kind: verifyPageInvalid})

		return
	}

	if _, err := h.authService.VerifyEmail(ctx, token); err != nil {
		if errors.Is(err, core_errors.ErrInvalidArgument) {
			log.Warn("verification token rejected", zap.Error(err))
			renderVerifyPage(rw, log, http.StatusBadRequest, verifyPageData{Kind: verifyPageInvalid})

			return
		}

		log.Error("failed to verify email", zap.Error(err))
		renderVerifyPage(rw, log, http.StatusInternalServerError, verifyPageData{Kind: verifyPageError})

		return
	}

	renderVerifyPage(rw, log, http.StatusOK, verifyPageData{Kind: verifyPageSuccess})
}
