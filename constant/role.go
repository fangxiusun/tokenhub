package constant

// User Roles
const (
	RoleGuestUser  = 0
	RoleCommonUser = 1
	RoleAdminUser  = 10
	RoleRootUser   = 100
)

// User Status
const (
	UserStatusEnabled  = 1
	UserStatusDisabled = 2
)

// IsValidateRole checks if the role is valid
func IsValidateRole(role int) bool {
	return role == RoleGuestUser || role == RoleCommonUser || role == RoleAdminUser || role == RoleRootUser
}

// CanManageRole checks if the given role can manage the target role
func CanManageRole(myRole, targetRole int) bool {
	return myRole == RoleRootUser || (myRole == RoleAdminUser && targetRole < RoleAdminUser)
}

