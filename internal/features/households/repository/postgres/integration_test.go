//go:build integration

package households_repository_postgres_test

import (
	core_domain "cohesive-core/internal/core/domain"
	households_repository_postgres "cohesive-core/internal/features/households/repository/postgres"
	"cohesive-core/internal/testhelpers"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateHousehold_CreatesOwnerMembership(t *testing.T) {
	pool := testhelpers.NewPostgres(t)
	repo := households_repository_postgres.NewHouseholdsRepository(pool)

	ownerID := testhelpers.InsertUser(t, pool, "owner@example.com")

	created, err := repo.CreateHousehold(
		context.Background(),
		core_domain.NewHouseholdUninitialized("Тестовый дом"),
		ownerID,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "Тестовый дом", created.Name)

	householdWithRole, err := repo.GetHouseholdByIDForUser(context.Background(), created.ID, ownerID)
	require.NoError(t, err)
	assert.Equal(t, core_domain.HouseholdRoleOwner, householdWithRole.Role)
}

func TestAcceptInvite_Success(t *testing.T) {
	pool := testhelpers.NewPostgres(t)
	repo := households_repository_postgres.NewHouseholdsRepository(pool)

	ownerID := testhelpers.InsertUser(t, pool, "owner@example.com")
	joinerID := testhelpers.InsertUser(t, pool, "joiner@example.com")

	household, err := repo.CreateHousehold(
		context.Background(),
		core_domain.NewHouseholdUninitialized("Дом"),
		ownerID,
	)
	require.NoError(t, err)

	invite, err := repo.CreateInvite(context.Background(), core_domain.NewHouseholdInvite(
		household.ID, "TESTCODE123", ownerID, time.Now().Add(time.Hour), nil,
	))
	require.NoError(t, err)

	joined, err := repo.AcceptInvite(context.Background(), invite.Code, joinerID)
	require.NoError(t, err)
	assert.Equal(t, household.ID, joined.ID)

	role, err := repo.GetMemberRole(context.Background(), household.ID, joinerID)
	require.NoError(t, err)
	assert.Equal(t, core_domain.HouseholdRoleMember, role)
}

func TestAcceptInvite_RejectsExpiredCode(t *testing.T) {
	pool := testhelpers.NewPostgres(t)
	repo := households_repository_postgres.NewHouseholdsRepository(pool)

	ownerID := testhelpers.InsertUser(t, pool, "owner@example.com")
	joinerID := testhelpers.InsertUser(t, pool, "joiner@example.com")

	household, err := repo.CreateHousehold(
		context.Background(),
		core_domain.NewHouseholdUninitialized("Дом"),
		ownerID,
	)
	require.NoError(t, err)

	invite, err := repo.CreateInvite(context.Background(), core_domain.NewHouseholdInvite(
		household.ID, "EXPIRED123", ownerID, time.Now().Add(-time.Hour), nil,
	))
	require.NoError(t, err)

	_, err = repo.AcceptInvite(context.Background(), invite.Code, joinerID)
	require.Error(t, err)
}

func TestAcceptInvite_RespectsMaxUses(t *testing.T) {
	pool := testhelpers.NewPostgres(t)
	repo := households_repository_postgres.NewHouseholdsRepository(pool)

	ownerID := testhelpers.InsertUser(t, pool, "owner@example.com")
	firstJoinerID := testhelpers.InsertUser(t, pool, "first@example.com")
	secondJoinerID := testhelpers.InsertUser(t, pool, "second@example.com")

	household, err := repo.CreateHousehold(
		context.Background(),
		core_domain.NewHouseholdUninitialized("Дом"),
		ownerID,
	)
	require.NoError(t, err)

	maxUses := 1
	invite, err := repo.CreateInvite(context.Background(), core_domain.NewHouseholdInvite(
		household.ID, "ONEUSE123", ownerID, time.Now().Add(time.Hour), &maxUses,
	))
	require.NoError(t, err)

	_, err = repo.AcceptInvite(context.Background(), invite.Code, firstJoinerID)
	require.NoError(t, err)

	_, err = repo.AcceptInvite(context.Background(), invite.Code, secondJoinerID)
	require.Error(t, err, "второй участник не должен пройти по уже исчерпанному инвайту")
}

func TestTransferOwnership_AtomicSwap(t *testing.T) {
	pool := testhelpers.NewPostgres(t)
	repo := households_repository_postgres.NewHouseholdsRepository(pool)

	ownerID := testhelpers.InsertUser(t, pool, "owner@example.com")
	memberID := testhelpers.InsertUser(t, pool, "member@example.com")

	household, err := repo.CreateHousehold(
		context.Background(),
		core_domain.NewHouseholdUninitialized("Дом"),
		ownerID,
	)
	require.NoError(t, err)

	invite, err := repo.CreateInvite(context.Background(), core_domain.NewHouseholdInvite(
		household.ID, "TRANSFER123", ownerID, time.Now().Add(time.Hour), nil,
	))
	require.NoError(t, err)

	_, err = repo.AcceptInvite(context.Background(), invite.Code, memberID)
	require.NoError(t, err)

	err = repo.TransferOwnership(context.Background(), household.ID, ownerID, memberID)
	require.NoError(t, err)

	newOwnerRole, err := repo.GetMemberRole(context.Background(), household.ID, memberID)
	require.NoError(t, err)
	assert.Equal(t, core_domain.HouseholdRoleOwner, newOwnerRole)

	oldOwnerRole, err := repo.GetMemberRole(context.Background(), household.ID, ownerID)
	require.NoError(t, err)
	assert.Equal(t, core_domain.HouseholdRoleAdmin, oldOwnerRole)
}

func TestTransferOwnership_RejectsWrongCaller(t *testing.T) {
	pool := testhelpers.NewPostgres(t)
	repo := households_repository_postgres.NewHouseholdsRepository(pool)

	ownerID := testhelpers.InsertUser(t, pool, "owner@example.com")
	notOwnerID := testhelpers.InsertUser(t, pool, "not-owner@example.com")
	targetID := testhelpers.InsertUser(t, pool, "target@example.com")

	household, err := repo.CreateHousehold(
		context.Background(),
		core_domain.NewHouseholdUninitialized("Дом"),
		ownerID,
	)
	require.NoError(t, err)

	err = repo.TransferOwnership(context.Background(), household.ID, notOwnerID, targetID)
	require.Error(t, err)
}
