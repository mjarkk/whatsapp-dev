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
	err = DB.Model(&models.Conversation{}).Preload("Messages.Buttons").First(&conversationResp, message.ConversationID).Error
	if err != nil {
		return err
	}

	webhook.NotivyMessage(message, false)

	return c.JSON(conversationResp)
}

func getConversationFromParam(c *fiber.Ctx) (*models.Conversation, error) {
	id, err := c.ParamsInt("id")
	if err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, errors.New("invalid conversation id")
	}

	conversation := &models.Conversation{}
	err = DB.Model(&models.Conversation{}).Preload("Messages.Buttons").First(conversation, id).Error
	return conversation, err
}

func CreateMessage(c *fiber.Ctx) error {
	request := struct {
		Message string
	}{}
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}
	if request.Message == "" {
		return errors.New("missing message")
	}

	conversation, err := getConversationFromParam(c)
	if err != nil {
		return err
	}

	newMessage := models.Message{
		ConversationID: conversation.ID,
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

func BtnQuickReply(c *fiber.Ctx) error {
	conversation, err := getConversationFromParam(c)
	if err != nil {
		return err
	}

	btnID, err := c.ParamsInt("btnId")
	if err != nil {
		return err
	}
	button := &models.MessageButton{}
	err = DB.Model(&models.MessageButton{}).First(button, btnID).Error
	if err != nil {
		return err
	}
	if button.ConversationID != conversation.ID {
		return errors.New("button does not belong to conversation")
	}

	newMessage := models.Message{
		ConversationID: conversation.ID,
		WhatsappID:     phonenumber.CreateWhatsappID(conversation.PhoneNumber),
		Direction:      models.DirectionOut,
		Message:        button.Text,
		Timestamp:      time.Now().Unix(),
		Payload:        button.Payload,
	}
	err = DB.Create(&newMessage).Error
	if err != nil {
		return err
	}

	conversation.Messages = append(conversation.Messages, newMessage)

	webhook.NotivyMessage(newMessage, false)

	return c.JSON(conversation)
}
