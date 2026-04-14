package channels

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
	group := app.Group("/channels")
	group.Use(middleware.Protected(), middleware.RequireSameTenant())

	group.Get("/", GetChannels)
	group.Post("/", CreateChannel)
	group.Put("/:id", UpdateChannel)
	group.Delete("/:id", DeleteChannel)
}

func GetChannels(c *fiber.Ctx) error {
	tenantIdHex := c.Locals("tenantId").(string)
	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := db.GetCollection("channels").Find(ctx, bson.M{"tenantId": tenantId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch channels"})
	}
	defer cursor.Close(ctx)

	var channels []models.Channel
	if err = cursor.All(ctx, &channels); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode channels"})
	}

	return c.JSON(channels)
}

func CreateChannel(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		RoomType string `json:"roomType"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	tenantIdHex := c.Locals("tenantId").(string)
	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	channel := models.Channel{
		ID:        primitive.NewObjectID(),
		TenantID:  tenantId,
		Name:      input.Name,
		Type:      input.Type,
		RoomType:  input.RoomType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.GetCollection("channels").InsertOne(ctx, channel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create channel"})
	}

	return c.Status(fiber.StatusCreated).JSON(channel)
}

func UpdateChannel(c *fiber.Ctx) error {
	channelIdHex := c.Params("id")
	channelId, err := primitive.ObjectIDFromHex(channelIdHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid channel ID"})
	}

	tenantIdHex := c.Locals("tenantId").(string)
	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	var input struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"name": input.Name, "updatedAt": time.Now()}}
	result, err := db.GetCollection("channels").UpdateOne(ctx, bson.M{"_id": channelId, "tenantId": tenantId}, update)

	if err != nil || result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Channel not found or unauthorized"})
	}

	return c.JSON(fiber.Map{"message": "Channel updated successfully"})
}

func DeleteChannel(c *fiber.Ctx) error {
	channelIdHex := c.Params("id")
	channelId, err := primitive.ObjectIDFromHex(channelIdHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid channel ID"})
	}

	tenantIdHex := c.Locals("tenantId").(string)
	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.GetCollection("channels").DeleteOne(ctx, bson.M{"_id": channelId, "tenantId": tenantId})
	if err != nil || result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Channel not found or unauthorized"})
	}

	return c.JSON(fiber.Map{"message": "Channel deleted successfully"})
}
