package webhooks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mjarkk/whatsapp-dev/go/db"
	"github.com/mjarkk/whatsapp-dev/go/lib/webhook"
	"github.com/mjarkk/whatsapp-dev/go/models"
)

func Test(c *fiber.Ctx) error {
	err := webhook.Validate()
	if err != nil {
		return err
	}

	message := models.Message{}
	err = db.DB.Model(models.Message{}).First(&message).Error
	if err != nil {
		// We can only test validating the webhook
		// There are no messages to send a webhook example message for
		return c.JSON(true)
	}

	err = webhook.NotivyMessage(message, true)
	if err != nil {
		return err
	}

	return c.JSON(true)
}
