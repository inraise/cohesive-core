package auth_transport_http

import (
	"time"

	core_domain "cohesive-core/internal/core/domain"

	"github.com/google/uuid"
)

type UserDTOResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Age       *int    `json:"age"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func userDTOFromDomain(user core_domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:        user.ID,
		Version:   user.Version,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
