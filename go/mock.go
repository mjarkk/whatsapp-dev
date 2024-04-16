package src

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mjarkk/whatsapp-dev/go/controller/messages"
)

func mockRoutes(r fiber.Router) {
	version := r.Group("/v:version")
	version.Post("/:phoneNumberId/messages", messages.Create)
}
