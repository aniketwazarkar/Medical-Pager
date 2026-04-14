package messages

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"medical-pager/internal/db"
	"medical-pager/internal/middleware"
	"medical-pager/internal/models"
	"medical-pager/internal/redis"
)

// RegisterRoutes links the fiber endpoints
func RegisterRoutes(app fiber.Router) {
	group := app.Group("/messages")
	// Must have JWT to access
	group.Use(middleware.Protected(), middleware.RequireSameTenant())

	group.Get("/:channelId", GetMessages)
	group.Post("/", SendMessage)
}

func GetMessages(c *fiber.Ctx) error {
	channelIdHex := c.Params("channelId")
	tenantIdHex := c.Locals("tenantId").(string)

	channelId, _ := primitive.ObjectIDFromHex(channelIdHex)
	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"channelId": channelId, "tenantId": tenantId}
	cursor, err := db.GetCollection("messages").Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}
	defer cursor.Close(ctx)

	var msgs []models.Message
	if err = cursor.All(ctx, &msgs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode messages"})
	}

	return c.JSON(msgs)
}

func SendMessage(c *fiber.Ctx) error {
	var input struct {
		ChannelID        string  `json:"channelId"`
		EncryptedContent string  `json:"encryptedContent"`
		MessageType      string  `json:"messageType"`
		Priority         string  `json:"priority"`
		PatientID        *string `json:"patientId,omitempty"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	tenantIdHex := c.Locals("tenantId").(string)
	senderIdHex := c.Locals("userId").(string)

	tenantId, _ := primitive.ObjectIDFromHex(tenantIdHex)
	senderId, _ := primitive.ObjectIDFromHex(senderIdHex)
	channelId, _ := primitive.ObjectIDFromHex(input.ChannelID)

	var patientId *primitive.ObjectID
	if input.PatientID != nil {
		pId, _ := primitive.ObjectIDFromHex(*input.PatientID)
		patientId = &pId
	}

	// Look up sender name so the message carries it for chat display
	var senderUser models.User
	senderCtx, senderCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer senderCancel()
	db.GetCollection("users").FindOne(senderCtx, bson.M{"_id": senderId}).Decode(&senderUser)

	// The content is already encrypted by the client, simply store it!
	encCont := input.EncryptedContent

	msg := models.Message{
		ID:               primitive.NewObjectID(),
		TenantID:         tenantId,
		ChannelID:        channelId,
		SenderID:         senderId,
		SenderName:       senderUser.Name,
		EncryptedContent: encCont,
		MessageType:      input.MessageType,
		Priority:         input.Priority,
		PatientID:        patientId,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.GetCollection("messages").InsertOne(ctx, msg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB insertion failed"})
	}

	// Publish to Redis Pub/Sub so WebSocket servers can dispatch to users
	payload, _ := json.Marshal(msg)
	redis.Client.Publish(redis.Ctx, "channel_"+channelId.Hex(), string(payload))

	return c.Status(fiber.StatusCreated).JSON(msg)
}
