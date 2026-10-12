package auth_service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"

	"go.uber.org/zap"
)

const verificationEmailSubject = "Подтверждение email в Cohesive"

func (s *AuthService) SendEmailVerification(
	ctx context.Context,
	user core_domain.User,
) error {
	plainToken, err := generateToken()
	if err != nil {
		return fmt.Errorf("generate verification token: %w", err)
	}

	expiresAt := time.Now().Add(s.verification.TokenTTL)

	link, err := buildVerificationLink(s.verification.BaseURL, plainToken)
	if err != nil {
		return fmt.Errorf("build verification link: %w", err)
	}

	verification := core_domain.NewEmailVerification(
		user.ID,
		user.Email,
		hashToken(plainToken),
		expiresAt,
	)

	if _, err = s.authRepository.CreateEmailVerification(ctx, verification); err != nil {
		return fmt.Errorf("save email verification: %w", err)
	}

	var (
		log  = core_logger.FromContextOrNop(ctx).With(zap.String("user_id", user.ID.String()))
		to   = user.Email
		body = buildVerificationBody(user.FirstName, link, expiresAt)
	)

	go func() {
		if err := s.mailer.Send(to, verificationEmailSubject, body); err != nil {
			log.Error("failed to send verification email", zap.Error(err))

			return
		}

		log.Debug("verification email sent")
	}()

	return nil
}

func (s *AuthService) VerifyEmail(
	ctx context.Context,
	tokenPlain string,
) (core_domain.User, error) {
	user, err := s.authRepository.VerifyEmailByHash(ctx, hashToken(tokenPlain))
	if err != nil {
		return core_domain.User{}, fmt.Errorf("verify email: %w", err)
	}

	return user, nil
}

func (s *AuthService) ResendVerification(
	ctx context.Context,
	email string,
) error {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("get user from repository: %w", err)
	}

	if user.IsVerified {
		return nil
	}

	if err = s.SendEmailVerification(ctx, user); err != nil {
		return fmt.Errorf("send email verification: %w", err)
	}

	return nil
}

func buildVerificationLink(baseURL, token string) (string, error) {
	link, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}

	query := link.Query()
	query.Set("token", token)
	link.RawQuery = query.Encode()

	return link.String(), nil
}

func buildVerificationBody(firstName, link string, expiresAt time.Time) string {
	return fmt.Sprintf(
		"Здравствуйте, %s!\r\n\r\n"+
			"Чтобы подтвердить email в Cohesive, перейдите по ссылке и нажмите кнопку «Подтвердить email»:\r\n%s\r\n\r\n"+
			"Ссылка действует до %s.\r\n\r\n"+
			"Если вы не регистрировались в Cohesive, просто проигнорируйте это письмо.",
		firstName,
		link,
		expiresAt.Format("02.01.2006 15:04 MST"),
	)
}
