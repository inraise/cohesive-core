package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SplitTypeExpenses string

const (
	Equal  SplitTypeExpenses = "equal"
	Custom SplitTypeExpenses = "custom"
)

type Expense struct {
	Id          uuid.UUID         `db:"id" json:"id"`
	FamilyID    uuid.UUID         `db:"family_id" json:"family_id"`
	Amount      decimal.Decimal   `db:"amount" json:"amount"`
	Currency    string            `db:"currency" json:"currency"`
	Category    string            `db:"category" json:"category"`
	Description string            `db:"description" json:"description"`
	Date        time.Time         `db:"date" json:"date"`
	PaidBy      uuid.UUID         `db:"paid_by" json:"paid_by"`
	SplitType   SplitTypeExpenses `db:"split_type" json:"split_type"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
	DeletedAt   time.Time         `db:"deleted_at" json:"deleted_at"`
}
