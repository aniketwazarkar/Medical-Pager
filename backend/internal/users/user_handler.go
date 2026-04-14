package users

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"medical-pager/internal/db"
	"medical-pager/internal/middleware"
	"medical-pager/internal/models"
)

func RegisterRoutes(app fiber.Router) {
	group := app.Group("/users")
	
	// Self-management (any authenticated user)
	group.Put("/me", middleware.Protected(), UpdateSelf)

	// Admin management (tenant_admin or super_admin only)
	adminGroup := group.Group("/")
	adminGroup.Use(middleware.Protected(), middleware.RequireSameTenant(), middleware.RequireRole(middleware.AdminRoles...))
	adminGroup.Get("/", GetUsers)
	adminGroup.Post("/", CreateUser)
	adminGroup.Put("/:id/role", UpdateUserRole)
}

func UpdateSelf(c *fiber.Ctx) error {
	userIdHex := c.Locals("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user context"})
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name cannot be empty"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"name": input.Name, "updatedAt": time.Now()}}
	result, err := db.GetCollection("users").UpdateOne(ctx, bson.M{"_id": userId}, update)

	if err != nil || result.MatchedCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	return c.JSON(fiber.Map{"message": "Profile updated successfully", "name": input.Name})
}

func CreateUser(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	validRoles := map[string]bool{
		middleware.RoleDoctor:      true,
		middleware.RoleNurse:       true,
		middleware.RoleStaff:       true,
		middleware.RoleTenantAdmin: true,
	}
	if !validRoles[input.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	tenantIdHex := c.Locals("tenantId").(string)
	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	user := models.User{
		ID:        primitive.NewObjectID(),
		TenantID:  tenantId,
		Email:     input.Email,
		Password:  string(hashed),
		Name:      input.Name,
		Role:      input.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.GetCollection("users").InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	user.Password = ""
	return c.Status(fiber.StatusCreated).JSON(user)
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

	// Validate allowed roles — super_admin cannot be assigned via this endpoint
	validRoles := map[string]bool{
		middleware.RoleDoctor:      true,
		middleware.RoleNurse:       true,
		middleware.RoleStaff:       true,
		middleware.RoleTenantAdmin: true,
	}
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
