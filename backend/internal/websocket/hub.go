package websocket

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"medical-pager/internal/redis"
)

// RegisterRoutes sets up the websocket endpoint
func RegisterRoutes(app fiber.Router) {
	// Middleware to ensure it's a websocket upgrade request
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			// In production, validate JWT here from headers or query param before upgrading
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:channelId", websocket.New(func(c *websocket.Conn) {
		channelId := c.Params("channelId")
		log.Printf("WS: ✅ Client connected to channel: %s", channelId)

		// Subscribe to Redis PubSub for this channel
		pubsub := redis.Client.Subscribe(redis.Ctx, "channel_"+channelId)
		defer func() {
			pubsub.Close()
			log.Printf("WS: ❌ Client disconnected from channel: %s", channelId)
		}()

		// Read loop for incoming WS messages
		go func() {
			for {
				var msg map[string]interface{}
				if err := c.ReadJSON(&msg); err != nil {
					log.Println("WS read error:", err)
					break
				}

				// Optional: publish received WS message to Redis
				// redis.Client.Publish(redis.Ctx, "channel_"+channelId, string_payload)
			}
		}()

		// Receive from Redis and send to WS Client
		ch := pubsub.Channel()
		for redisMsg := range ch {
			log.Printf("WS: 🚀 Pushing message onto channel %s -> Payload: %s", channelId, redisMsg.Payload)
			if err := c.WriteMessage(websocket.TextMessage, []byte(redisMsg.Payload)); err != nil {
				log.Println("WS write error:", err)
				break
			}
		}
	}))
}
