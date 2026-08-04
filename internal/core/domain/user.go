package core_domain

import (
	core_errors "cohesive-core/internal/core/errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID      uuid.UUID
	Version int

	Email        string
	PasswordHash string

	FirstName string
	LastName  *string
	Age       *int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(
	ID uuid.UUID,
	version int,
	email, passwordHash, firstName string,
	lastName *string,
	age *int,
	createdAt, updatedAt time.Time,
) User {
	return User{
		ID:           ID,
		Version:      version,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Age:          age,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func NewUserUninitialized(email, passwordHash, firstName string, lastName *string, age *int) User {
	return User{
		ID:           UninitializedID,
		Version:      UninitializedVersion,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Age:          age,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (u *User) Validate() error {
	emailLen := len([]rune(u.Email))
	if emailLen < 5 && emailLen > 100 {
		return fmt.Errorf(
			"invalid `Email` len: %d: %w",
			emailLen,
			core_errors.ErrInvalidArgument,
		)
	}

	firstNameLen := len([]rune(u.FirstName))
	if firstNameLen < 5 && firstNameLen > 100 {
		return fmt.Errorf(
			"invalid `FirstName` len: %d: %w",
			firstNameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.LastName != nil {
		lastNameLen := len([]rune(*u.LastName))
		if lastNameLen < 5 && lastNameLen > 100 {
			return fmt.Errorf(
				"invalid `LastName` len: %d: %w",
				lastNameLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if u.UpdatedAt.Before(u.CreatedAt) {
		return fmt.Errorf(
			"`UpdatedAt` is before `CreatedAt`: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
