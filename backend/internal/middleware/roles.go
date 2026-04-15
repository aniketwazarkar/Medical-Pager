package middleware

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
)

// -----------------------------------------------------------------
// Embed the authoritative role definitions at compile time.
// File: backend/internal/middleware/roles.json
// -----------------------------------------------------------------

//go:embed roles.json
var rolesJSON []byte

// RoleDefinition mirrors one entry from data/roles.json.
type RoleDefinition struct {
	Value       string   `json:"value"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Assignable  bool     `json:"assignable"`
	Groups      []string `json:"groups"`
}

// Roles holds all role definitions loaded from JSON.
var Roles []RoleDefinition

// AssignableRoles is a fast-lookup set of roles that can be assigned to users
// (i.e. excludes super_admin which is never assignable via normal endpoints).
var AssignableRoles map[string]bool

// Role constants — kept for type-safe use across the codebase.
// These are validated against the JSON at init so they stay in sync.
const (
	RoleSuperAdmin  = "super_admin"
	RoleTenantAdmin = "tenant_admin"
	RoleDoctor      = "doctor"
	RoleNurse       = "nurse"
	RoleStaff       = "staff"
)

// Derived role group slices — built from the JSON at init.
var (
	AllRoles      []string
	AdminRoles    []string
	SuperOnly     []string
	ClinicalRoles []string
)

func init() {
	if err := json.Unmarshal(rolesJSON, &Roles); err != nil {
		log.Fatalf("[roles] Failed to parse data/roles.json: %v", err)
	}

	AssignableRoles = make(map[string]bool, len(Roles))

	groupSet := map[string]*[]string{
		"AllRoles":      &AllRoles,
		"AdminRoles":    &AdminRoles,
		"SuperOnly":     &SuperOnly,
		"ClinicalRoles": &ClinicalRoles,
	}

	for _, r := range Roles {
		if r.Assignable {
			AssignableRoles[r.Value] = true
		}
		for _, g := range r.Groups {
			if slice, ok := groupSet[g]; ok {
				*slice = append(*slice, r.Value)
			}
		}
	}

	log.Printf("[roles] Loaded %d roles from data/roles.json (%d assignable)",
		len(Roles), len(AssignableRoles))
}

// IsValidRole reports whether role is a recognised system role.
func IsValidRole(role string) bool {
	for _, r := range Roles {
		if r.Value == role {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------
// GET /api/v1/roles
// Returns the list of roles that are assignable by admins.
// Requires a valid JWT (Protected middleware applied in routes.go).
// -----------------------------------------------------------------

func GetAssignableRoles(c *fiber.Ctx) error {
	var assignable []RoleDefinition
	for _, r := range Roles {
		if r.Assignable {
			assignable = append(assignable, r)
		}
	}
	return c.JSON(assignable)
}
