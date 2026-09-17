package budget_repository_postgres

import (
	"fmt"

	core_domain "cohesive-core/internal/core/domain"
)

type scanner interface {
	Scan(dest ...any) error
}

func scanTransaction(row scanner) (core_domain.HouseholdTransaction, error) {
	var model TransactionModel

	err := row.Scan(
		&model.ID,
		&model.HouseholdID,
		&model.Version,
		&model.Type,
		&model.Amount,
		&model.Description,
		&model.CreatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return core_domain.HouseholdTransaction{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.NewHouseholdTransaction(
		model.ID,
		model.HouseholdID,
		model.Version,
		core_domain.TransactionType(model.Type),
		model.Amount,
		model.Description,
		model.CreatedBy,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}
