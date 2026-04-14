package encryption

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"medical-pager/internal/db"
	"medical-pager/internal/middleware"
	"medical-pager/internal/models"
)

func RegisterRoutes(app fiber.Router) {
	group := app.Group("/admin")
	// Require JWT and SuperAdmin or TenantAdmin
	group.Use(middleware.Protected(), middleware.RequireRole("super_admin", "tenant_admin"))
	group.Post("/decrypt-message", DecryptMessage)
}

func DecryptMessage(c *fiber.Ctx) error {
	var input struct {
		MessageID string `json:"messageId"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userIdHex := c.Locals("userId").(string)
	tenantIdHex := c.Locals("tenantId").(string)

	uId, _ := primitive.ObjectIDFromHex(userIdHex)
	tId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Audit Log
	db.GetCollection("audit_logs").InsertOne(ctx, models.AuditLog{
		TenantID:  tId,
		UserID:    uId,
		Action:    "DECRYPT_MESSAGE",
		CreatedAt: time.Now(),
		Metadata: fiber.Map{
			"messageId": input.MessageID,
		},
	})

	// Fetch message
	_ , _ = primitive.ObjectIDFromHex(input.MessageID)
	// Add tenant barrier: role=tenant_admin can only decrypt their own tenant
	// If super_admin, could bypass, but keeping simple for now.

	return c.JSON(fiber.Map{
		"messageType": "example decrypt stub",
		"status":      "In a real app, query message, call encryption.Decrypt, and return plaintext",
	})
}
