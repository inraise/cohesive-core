package households_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// GetHouseholdByIDForUser отдаёт дом только если userID реально состоит
// в нём. Если дома не существует ИЛИ пользователь не его участник - в обоих
// случаях возвращается один и тот же ErrNotFound, чтобы нельзя было
// перебором id узнать о существовании чужих домов.
func (r *HouseholdsRepository) GetHouseholdByIDForUser(
	ctx context.Context,
	householdID uuid.UUID,
	userID uuid.UUID,
) (core_domain.HouseholdWithRole, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT h.id, h.version, h.name, h.created_at, h.updated_at, hm.role
		FROM households h
		JOIN household_members hm ON hm.household_id = h.id
		WHERE h.id = $1 AND hm.user_id = $2;
	`

	row := r.pool.QueryRow(ctx, query, householdID, userID)

	var (
		householdModel HouseholdModel
		role           string
	)

	err := row.Scan(
		&householdModel.ID,
		&householdModel.Version,
		&householdModel.Name,
		&householdModel.CreatedAt,
		&householdModel.UpdatedAt,
		&role,
	)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.HouseholdWithRole{}, fmt.Errorf(
				"household with id %q for user %q: %w",
				householdID, userID, core_errors.ErrNotFound,
			)
		}

		return core_domain.HouseholdWithRole{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.HouseholdWithRole{
		Household: core_domain.NewHousehold(
			householdModel.ID,
			householdModel.Version,
			householdModel.Name,
			householdModel.CreatedAt,
			householdModel.UpdatedAt,
		),
		Role: core_domain.HouseholdRole(role),
	}, nil
}
