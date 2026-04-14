package users

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"medical-pager/internal/db"
	"medical-pager/internal/middleware"
	"medical-pager/internal/models"
)

func RegisterRoutes(app fiber.Router) {
	group := app.Group("/users")
	
	// Protected by JWT and requires 'tenant_admin' or 'super_admin' roles
	group.Use(middleware.Protected(), middleware.RequireSameTenant(), middleware.RequireRole("tenant_admin", "super_admin"))

	group.Get("/", GetUsers)
	group.Put("/:id/role", UpdateUserRole)
}

func GetUsers(c *fiber.Ctx) error {
	tenantIdHex := c.Locals("tenantId").(string)
	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.GetCollection("users").Find(ctx, bson.M{"tenantId": tenantId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode users"})
	}

	// We apply primitive masking since Find opts wasn't strictly typed above
	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(users)
}

func UpdateUserRole(c *fiber.Ctx) error {
	userIdHex := c.Params("id")
	userId, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	tenantIdHex := c.Locals("tenantId").(string)
	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	var input struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Validate allowed roles
	validRoles := map[string]bool{"doctor": true, "nurse": true, "staff": true, "tenant_admin": true}
	if !validRoles[input.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role assignment"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"role": input.Role, "updatedAt": time.Now()}}
	result, err := db.GetCollection("users").UpdateOne(ctx, bson.M{"_id": userId, "tenantId": tenantId}, update)

	if err != nil || result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found or unauthorized"})
	}

	return c.JSON(fiber.Map{"message": "User role updated successfully"})
}
