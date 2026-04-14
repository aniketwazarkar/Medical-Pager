package api

import (
	"github.com/gofiber/fiber/v2"

	"medical-pager/internal/auth"
	"medical-pager/internal/channels"
	"medical-pager/internal/encryption"
	"medical-pager/internal/messages"
	"medical-pager/internal/tenants"
	"medical-pager/internal/users"
	"medical-pager/internal/websocket"
)

// SetupRoutes registers all app routes
func SetupRoutes(app *fiber.App) {
	// API Group
	api := app.Group("/api/v1")

	auth.RegisterRoutes(api)
	messages.RegisterRoutes(api)
	encryption.RegisterRoutes(api)
	websocket.RegisterRoutes(api)

	// New dynamic bindings
	channels.RegisterRoutes(api)
	users.RegisterRoutes(api)
	tenants.RegisterRoutes(api)
}
