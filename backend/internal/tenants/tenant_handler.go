package tenants

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"medical-pager/internal/db"
	"medical-pager/internal/middleware"
	"medical-pager/internal/models"
)

func RegisterRoutes(app fiber.Router) {
	group := app.Group("/tenants")
	
	// Protected by JWT and requires strict 'super_admin' role
	group.Use(middleware.Protected(), middleware.RequireRole(middleware.SuperOnly...))

	group.Get("/", GetTenants)
}

func GetTenants(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.GetCollection("tenants").Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tenants"})
	}
	defer cursor.Close(ctx)

	var tenants []models.Tenant
	if err = cursor.All(ctx, &tenants); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode tenants"})
	}

	return c.JSON(tenants)
}
