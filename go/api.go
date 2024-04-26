package src

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mjarkk/whatsapp-dev/go/controller/conversations"
	"github.com/mjarkk/whatsapp-dev/go/controller/templates"
	"github.com/mjarkk/whatsapp-dev/go/controller/webhooks"
	"github.com/mjarkk/whatsapp-dev/go/controller/websocket"
	"github.com/mjarkk/whatsapp-dev/go/state"
)

func apiRoutes(r fiber.Router) {
	r.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}

		return c.Status(400).JSON(ErrorResponse{
			Error: err.Error(),
		})
	})

	r.Get("/events", websocket.EventsRoute)

	r.Get("/info", func(c *fiber.Ctx) error {
		return c.JSON(struct {
			GraphToken         string `json:"graphToken"`
			AppSecret          string `json:"appSecret"`
			PhoneNumber        string `json:"phoneNumber"`
			PhoneNumberID      string `json:"phoneNumberID"`
			WebhookURL         string `json:"webhookURL"`
			WebhookVerifyToken string `json:"webhookVerifyToken"`
		}{
			GraphToken:         state.GraphToken.Get(),
			AppSecret:          state.AppSecret.Get(),
			PhoneNumber:        state.PhoneNumber.Get(),
			PhoneNumberID:      state.PhoneNumberID.Get(),
			WebhookURL:         state.WebhookURL.Get(),
			WebhookVerifyToken: state.WebhookVerifyToken.Get(),
		})
	})

	r.Get("/conversations", conversations.Index)
	r.Post("/conversations", conversations.Create)
	r.Post("/conversations/:id", conversations.CreateMessage)
	r.Post("/conversations/:id/btnQuickReply/:btnId", conversations.BtnQuickReply)

	r.Get("/templates", templates.Index)
	r.Post("/templates", templates.Create)
	r.Patch("/templates/:id", templates.Update)
	r.Delete("/templates/:id", templates.Delete)

	r.Post("/webhook/test", webhooks.Test)
}
