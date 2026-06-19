package models

import (
	"time"

	"github.com/google/uuid"
)

type RoleMembers string

const (
	Admin  RoleMembers = "admin"
	Member RoleMembers = "member"
	Child  RoleMembers = "child"
)

type FamilyMember struct {
	Id          uuid.UUID    `db:"id" json:"id"`
	FamilyID    uuid.UUID    `db:"family_id" json:"family_id"`
	UserId       uuid.UUID       `db:"user_id" json:"user_id"`
	Role RoleMembers       `db:"role" json:"role"`
	JoinedAt   time.Time    `db:"joined_at" json:"joined_at"`
}
