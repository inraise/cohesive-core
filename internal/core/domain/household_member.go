package core_domain

type HouseholdRole string

const (
	HouseholdRoleOwner  HouseholdRole = "owner"
	HouseholdRoleAdmin  HouseholdRole = "admin"
	HouseholdRoleMember HouseholdRole = "member"
)

type HouseholdWithRole struct {
	Household Household
	Role      HouseholdRole
}
