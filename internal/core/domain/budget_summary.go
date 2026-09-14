package core_domain

import "github.com/google/uuid"

type MemberContribution struct {
	UserID uuid.UUID
	Amount int64
}
