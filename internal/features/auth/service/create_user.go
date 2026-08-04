package auth_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Email     string  `json:"email" validate:"required,min=5,max=100"`
	Password  string  `json:"password" validate:"required,min=10,max=100"`
	FirstName string  `json:"first_name" validate:"required,min=3,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,min=3,max=100"`
	Age       *int    `json:"age" validate:"omitempty,min=0,max=130"`
}

func (s *AuthService) CreateUser(
	ctx context.Context,
	user CreateUserRequest,
) (core_domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("hash user password: %w", err)
	}

	userUninit := user
	userUninit.Password = string(hashedPassword)
	userDomain := domainFromDTO(userUninit)

	if err := userDomain.Validate(); err != nil {
		return core_domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	userResp, err := s.authRepository.CreateUser(ctx, userDomain)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return userResp, nil
}

func domainFromDTO(dto CreateUserRequest) core_domain.User {
	return core_domain.NewUserUninitialized(
		dto.Email,
		dto.Password,
		dto.FirstName,
		dto.LastName,
		dto.Age,
	)
}
