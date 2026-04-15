package tenants

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"medical-pager/internal/db"
	"medical-pager/internal/middleware"
	"medical-pager/internal/models"
	
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(app fiber.Router) {
	group := app.Group("/tenants")
	
	// Protected by JWT and requires strict 'super_admin' role
	group.Use(middleware.Protected(), middleware.RequireRole(middleware.SuperOnly...))

	group.Get("/", GetTenants)
	group.Post("/", CreateTenant)
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

func CreateTenant(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req struct {
		Name          string `json:"name"`
		Domain        string `json:"domain"`
		AdminName     string `json:"adminName"`
		AdminEmail    string `json:"adminEmail"`
		AdminPassword string `json:"adminPassword"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" || req.AdminEmail == "" || req.AdminPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	newTenant := models.Tenant{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Domain:    req.Domain,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := db.GetCollection("tenants").InsertOne(ctx, newTenant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create tenant"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash admin password"})
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		TenantID:  newTenant.ID,
		Name:      req.AdminName,
		Email:     req.AdminEmail,
		Password:  string(hashedPassword),
		Role:      "tenant_admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = db.GetCollection("users").InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Tenant created, but failed to create admin user"})
	}

	return c.Status(fiber.StatusCreated).JSON(newTenant)
}
