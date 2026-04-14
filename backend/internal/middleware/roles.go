package middleware

// Role constants — single source of truth for all roles in the system.
// All authorization decisions are enforced here on the backend;
// the frontend only reads these values from the JWT for UX purposes.
const (
	RoleSuperAdmin   = "super_admin"
	RoleTenantAdmin  = "tenant_admin"
	RoleDoctor       = "doctor"
	RoleNurse        = "nurse"
	RoleStaff        = "staff"
)

// AllRoles lists every valid role in the system.
var AllRoles = []string{
	RoleSuperAdmin,
	RoleTenantAdmin,
	RoleDoctor,
	RoleNurse,
	RoleStaff,
}

// AdminRoles are roles that can manage tenants and users.
var AdminRoles = []string{RoleSuperAdmin, RoleTenantAdmin}

// SuperOnly restricts to system owner only.
var SuperOnly = []string{RoleSuperAdmin}

// ClinicalRoles are roles that can access clinical channels and messages.
var ClinicalRoles = []string{RoleSuperAdmin, RoleTenantAdmin, RoleDoctor, RoleNurse, RoleStaff}

// IsValidRole checks whether a role string is a recognized system role.
func IsValidRole(role string) bool {
	for _, r := range AllRoles {
		if r == role {
			return true
		}
	}
	return false
}
