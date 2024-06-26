package models

import (
	"errors"

	. "github.com/mjarkk/whatsapp-dev/go/db"
	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	PhoneNumberId string    `json:"phoneNumberId"`
	PhoneNumber   string    `json:"phoneNumber"`
	Messages      []Message `json:"messages"`
}

type Message struct {
	gorm.Model
	ConversationID uint            `json:"conversationId"`
	WhatsappID     string          `json:"whatsappID"`
	Direction      Direction       `json:"direction"`
	HeaderMessage  *string         `json:"headerMessage"`
	Message        string          `json:"message"`
	FooterMessage  *string         `json:"footerMessage"`
	Timestamp      int64           `json:"timestamp"`
	Payload        *string         `json:"payload"`
	Buttons        []MessageButton `json:"buttons"`
}

type MessageButton struct {
	gorm.Model
	ConversationID uint    `json:"conversationId"`
	MessageID      uint    `json:"messageId"`
	Text           string  `json:"text"`
	Payload        *string `json:"payload"`
}

type Direction string

const (
	DirectionIn  Direction = "in"
	DirectionOut Direction = "out"
)

func (m *Message) CreateOrAppend(number string) error {
	conversationID := uint(0)

	exsistingConversation := Conversation{}
	err := DB.Model(&Conversation{}).Where("phone_number = ?", number).First(&exsistingConversation).Error
	if err == nil {
		// Append messsage to exsisting conversation
		conversationID = exsistingConversation.ID
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Start a new converstaion
		newConversation := Conversation{
			PhoneNumberId: number,
			PhoneNumber:   number,
		}
		err = DB.Create(&newConversation).Error
		if err != nil {
			return err
		}

		conversationID = newConversation.ID
	} else {
		return err
	}

	m.ConversationID = conversationID

	for idx, button := range m.Buttons {
		button.ConversationID = conversationID
		m.Buttons[idx] = button
	}

	err = DB.Create(m).Error
	if err != nil {
		return err
	}

	return nil
}
