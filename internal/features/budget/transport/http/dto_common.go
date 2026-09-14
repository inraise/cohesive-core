package budget_transport_http

import (
	"time"

	core_domain "cohesive-core/internal/core/domain"

	"github.com/google/uuid"
)

type TransactionDTOResponse struct {
	ID          uuid.UUID `json:"id"`
	HouseholdID uuid.UUID `json:"household_id"`
	Version     int       `json:"version"`

	Type        string    `json:"type"`
	Amount      int64     `json:"amount"`
	Description *string   `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func transactionDTOFromDomain(tx core_domain.HouseholdTransaction) TransactionDTOResponse {
	return TransactionDTOResponse{
		ID:          tx.ID,
		HouseholdID: tx.HouseholdID,
		Version:     tx.Version,
		Type:        string(tx.Type),
		Amount:      tx.Amount,
		Description: tx.Description,
		CreatedBy:   tx.CreatedBy,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}
}

type MemberContributionDTOResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Amount int64     `json:"amount"`
}

type BudgetSummaryDTOResponse struct {
	Balance       int64                           `json:"balance"`
	TotalDeposits int64                           `json:"total_deposits"`
	TotalExpenses int64                           `json:"total_expenses"`
	Contributions []MemberContributionDTOResponse `json:"contributions"`
}
