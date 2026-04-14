package auth

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"medical-pager/internal/db"
	"medical-pager/internal/middleware"
	"medical-pager/internal/models"
	"medical-pager/utils"
)

// RegisterRoutes sets up the auth routes
func RegisterRoutes(app fiber.Router) {
	group := app.Group("/auth")
	group.Post("/login", Login)
	group.Post("/register", Register)
	// group.Post("/refresh", Refresh)
}

func generateToken(user *models.User) (string, error) {
	secret := utils.GetEnv("JWT_SECRET", "supersecretjwtkey_change_in_production")
	claims := middleware.Claims{
		UserID:   user.ID.Hex(),
		TenantID: user.TenantID.Hex(),
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := db.GetCollection("users").FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := generateToken(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	// Audit Log
	db.GetCollection("audit_logs").InsertOne(ctx, models.AuditLog{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		Action:    "LOGIN",
		CreatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"ip":        c.IP(),
			"userAgent": string(c.Request().Header.UserAgent()),
		},
	})

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func Register(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		TenantID string `json:"tenantId"`
		Role     string `json:"role"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	tId, _ := primitive.ObjectIDFromHex(input.TenantID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Optionally check if tenant exists
	// var tenant models.Tenant
	// err = db.GetCollection("tenants").FindOne(ctx, bson.M{"_id": tId}).Decode(&tenant)
	// if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tenant does not exist"}) }

	user := models.User{
		ID:        primitive.NewObjectID(),
		TenantID:  tId,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Name:      input.Name,
		Role:      input.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Make sure email is unique via index or direct check
	opts := options.Update().SetUpsert(true)
	_, err = db.GetCollection("users").UpdateOne(ctx, bson.M{"email": input.Email}, bson.M{"$setOnInsert": user}, opts)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
