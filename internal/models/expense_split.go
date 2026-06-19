package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpenseSplits struct {
	Id         uuid.UUID       `db:"id" json:"id"`
	ExpenseId  uuid.UUID       `db:"expense_id" json:"expense_id"`
	UserId     uuid.UUID       `db:"user_id" json:"user_id"`
	Percentage decimal.Decimal `db:"percentage" json:"percentage"`
}
