package conversations

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	. "github.com/mjarkk/whatsapp-dev/go/db"
	"github.com/mjarkk/whatsapp-dev/go/lib/webhook"
	"github.com/mjarkk/whatsapp-dev/go/models"
	"github.com/mjarkk/whatsapp-dev/go/utils/phonenumber"
)

func Index(c *fiber.Ctx) error {
	conversation := []models.Conversation{}
	err := DB.Model(&models.Conversation{}).Preload("Messages.Buttons").Find(&conversation).Error
	if err != nil {
		return err
	}

	return c.JSON(conversation)
}

func Create(c *fiber.Ctx) error {
	request := struct {
		PhoneNumber string `json:"phoneNumber"`
		Message     string `json:"message"`
	}{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.PhoneNumber == "" {
		return errors.New("missing phone number")
	}
	if request.Message == "" {
		return errors.New("missing message")
	}

	parsedPhoneNumber, err := phonenumber.Parse(request.PhoneNumber, true)
	if err != nil {
		return err
	}

	message := models.Message{
		WhatsappID: parsedPhoneNumber.WhatsappMessageID,
		Direction:  models.DirectionOut,
		Message:    request.Message,
		Timestamp:  time.Now().Unix(),
	}
	err = message.CreateOrAppend(parsedPhoneNumber.Parsed)
	if err != nil {
		return err
	}

	conversationResp := models.Conversation{}
	err = DB.Model(&models.Conversation{}).Preload("Messages").First(&conversationResp, message.ConversationID).Error
	if err != nil {
		return err
	}

	webhook.NotivyMessage(message, false)

	return c.JSON(conversationResp)
}

func CreateMessage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id <= 0 {
		return errors.New("invalid conversation id")
	}

	request := struct {
		Message string
	}{}
	err = c.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Message == "" {
		return errors.New("missing message")
	}

	conversation := models.Conversation{}
	err = DB.Model(&models.Conversation{}).Preload("Messages").First(&conversation, id).Error
	if err != nil {
		return err
	}

	newMessage := models.Message{
		ConversationID: uint(id),
		WhatsappID:     phonenumber.CreateWhatsappID(conversation.PhoneNumber),
		Direction:      models.DirectionOut,
		Message:        request.Message,
		Timestamp:      time.Now().Unix(),
	}

	err = DB.Create(&newMessage).Error
	if err != nil {
		return err
	}

	conversation.Messages = append(conversation.Messages, newMessage)

	webhook.NotivyMessage(newMessage, false)

	return c.JSON(conversation)
}
