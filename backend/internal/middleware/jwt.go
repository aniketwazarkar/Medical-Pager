package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"medical-pager/utils"
)

// Claims represents the JWT claims for the application
type Claims struct {
	UserID   string `json:"userId"`
	TenantID string `json:"tenantId"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Protected requires a valid JWT token
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}

		tokenString := strings.Split(authHeader, " ")[1]
		secret := utils.GetEnv("JWT_SECRET", "supersecretjwtkey_change_in_production")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Set context locals for next handlers
		c.Locals("userId", claims.UserID)
		c.Locals("tenantId", claims.TenantID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RequireRole enforces specific RBAC roles
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}
}

// RequireSameTenant ensures that a user can only query their own tenant boundary if not super_admin
func RequireSameTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, _ := c.Locals("role").(string)
		if userRole == "super_admin" {
			return c.Next()
		}

		queryTenantId := c.Query("tenantId")
		userTenantId, _ := c.Locals("tenantId").(string)

		// Note: Most routes will implicitly use the userTenantId for DB filtering to enforce isolation.
		// If explicitly passing tenantId, validate it matches.
		if queryTenantId != "" && queryTenantId != userTenantId {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Cross-tenant access forbidden"})
		}

		return c.Next()
	}
}
