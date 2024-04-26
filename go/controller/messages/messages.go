package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mjarkk/whatsapp-dev/go/controller/websocket"
	. "github.com/mjarkk/whatsapp-dev/go/db"
	"github.com/mjarkk/whatsapp-dev/go/models"
	"github.com/mjarkk/whatsapp-dev/go/state"
	"github.com/mjarkk/whatsapp-dev/go/utils/phonenumber"
)

func parseVersion(c *fiber.Ctx) error {
	versionParam := c.Params("version")
	versionParts := strings.Split(versionParam, ".")
	majorVersion, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return err
	}
	if majorVersion <= 10 {
		return errors.New("incorrect facebook grapth api version, only major version higher than 10 are supported")
	}
	return nil
}

type ErrorKind uint8

const (
	AuthTokenMalformed ErrorKind = iota
	AuthTokenInvalidAuthKind
	AuthTokenMissingAuthKind
	AuthTokenCannotBeDecrypted
	AuthTokenInvalidContentType
	RecipientPhoneNumberNotAllowed
)

func ErrValues(kind ErrorKind) (status int, code int, message string) {
	switch kind {
	case AuthTokenMalformed:
		return 400, 190, "Malformed access token"
	case AuthTokenInvalidAuthKind:
		return 401, 190, "Invalid auth type in access token"
	case AuthTokenMissingAuthKind:
		return 400, 190, "Missing authentication header"
	case AuthTokenCannotBeDecrypted:
		return 401, 190, "The access token could not be decrypted"
	case AuthTokenInvalidContentType:
		return 400, 190, "Invalid content type (application/json)"
	case RecipientPhoneNumberNotAllowed:
		return 400, 131030, "(#131030) Recipient phone number not in allowed list"
	default:
		return 400, 102, "Unknown error kind"
	}
}

func authError(c *fiber.Ctx, kind ErrorKind) error {
	httpStatusCode, code, message := ErrValues(kind)

	c.Response().Header.Set("www-authenticate", `OAuth "Facebook Platform" "invalid_request" "`+message+`"`)

	return c.Status(httpStatusCode).JSON(map[string]any{
		"error": map[string]any{
			"message":    message,
			"type":       "OAuthException",
			"code":       code,
			"fbtrace_id": "MDAwMDAwMDAwMDAwMDAwMDAw",
		}})
}

func customError(c *fiber.Ctx, message string, details ...string) error {
	c.Response().Header.Set("www-authenticate", `OAuth "Facebook Platform" "invalid_request" "`+message+`"`)

	errData := map[string]any{
		"message":    message,
		"type":       "OAuthException",
		"code":       100,
		"fbtrace_id": "MDAwMDAwMDAwMDAwMDAwMDAw",
	}

	if len(details) > 0 {
		errData["error_data"] = map[string]any{
			"messaging_product": "whatsapp",
			"details":           details[0],
		}
	}

	return c.Status(400).JSON(map[string]any{"error": errData})
}

func Create(c *fiber.Ctx) error {
	err := parseVersion(c)
	if err != nil {
		return err
	}

	c.Response().Header.Set("facebook-api-version", "v18.0")

	headers := c.GetReqHeaders()
	authHeader := ""
	hasContentType := false
	for key, value := range headers {
		switch strings.ToLower(key) {
		case "content-type":
			if strings.ToLower(value) != "application/json" {
				return authError(c, AuthTokenInvalidContentType)
			}
			hasContentType = true
		case "authorization":
			authHeader = value
		}
	}
	if !hasContentType {
		return authError(c, AuthTokenInvalidContentType)
	}

	// Validate auth token
	if authHeader == "" {
		return authError(c, AuthTokenMissingAuthKind)
	}
	_, after, found := strings.Cut(authHeader, "Bearer ")
	if !found {
		return authError(c, AuthTokenInvalidAuthKind)
	}
	if state.GraphToken.Get() != after {
		return authError(c, AuthTokenMalformed)
	}

	// Validate request content

	bodyBytes := c.Body()
	body := struct {
		MessagingProduct string           `json:"messaging_product"`
		To               string           `json:"to"`
		Type             string           `json:"type"` // "template", "text"
		Template         *TemplateOptions `json:"template"`
		Text             *TextOptions     `json:"text"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return customError(c, "(#100) The parameter messaging_product is required.", "Invalid JSON, err: "+err.Error())
	}
	if body.To == "" {
		return customError(c, "The parameter to is required.")
	}
	if strings.ToLower(body.MessagingProduct) != "whatsapp" {
		messagingProductJSON, _ := json.Marshal(body.MessagingProduct)
		messagingProductJSONStr := string(messagingProductJSON)
		errMsg := fmt.Sprintf("(#100) Param messaging_product must be one of {WHATSAPP} - got %s.", messagingProductJSONStr)
		return customError(c, errMsg)
	}

	to, err := phonenumber.Parse(body.To, false)
	if err != nil {
		return authError(c, RecipientPhoneNumberNotAllowed)
	}

	switch strings.ToLower(body.Type) {
	case "", "text":
		if body.Text == nil {
			return customError(c, "(#100) Invalid parameter", "Parameter 'text' is mandatory for type 'text'")
		}
		return handleSendTextMessage(c, *body.Text, to)
	case "template":
		if body.Template == nil {
			return customError(c, "(#100) Invalid parameter", "Parameter 'template' is mandatory for type 'template'")
		}
		return handleSendTemplateMessage(c, *body.Template, to)
	default:
		return customError(c, "(#100) Invalid parameter", "Parameter 'type' must be one of {TEXT, TEMPLATE}")
	}
}

type TextOptions struct {
	Body string `json:"body"`
}

func handleSendTextMessage(c *fiber.Ctx, text TextOptions, to *phonenumber.ParsedPhoneNumber) error {
	if text.Body == "" {
		return customError(c, "(#100) The parameter text['body'] is required.")
	}

	conversation := models.Conversation{}
	err := DB.Model(&models.Conversation{}).First(&conversation, "phone_number = ?", to.Parsed).Error
	if err != nil {
		return authError(c, RecipientPhoneNumberNotAllowed)
	}

	message := &models.Message{
		ConversationID: uint(conversation.ID),
		WhatsappID:     to.WhatsappMessageID,
		Direction:      models.DirectionIn,
		Message:        text.Body,
		Timestamp:      time.Now().Unix(),
	}
	err = DB.Create(message).Error
	if err != nil {
		return customError(c, "(#100) WhatsApp-Dev Error creating message", err.Error())
	}

	websocket.SendMessage(*message)

	// FIXME notify webhook

	return c.JSON(map[string]any{
		"messaging_product": "whatsapp",
		"contacts": []map[string]string{{
			"input": to.Original,
			"wa_id": to.Parsed,
		}},
		"messages": []map[string]string{{
			"id":             message.WhatsappID,
			"message_status": "accepted",
		}},
	})
}

type TemplateOptions struct {
	Name     string `json:"name"` // "hello_world",
	Language struct {
		Code   string `json:"code"`   // "en_US", "en"
		Policy string `json:"policy"` // "deterministic"
	} `json:"language"` // { "code": "en_US" }
	Components []TemplateComponent `json:"components"`
}

type TemplateComponent struct {
	Type       string `json:"type"`     // "body", "button"
	SubType    string `json:"sub_type"` // "quick_reply" (in case of button)
	Index      string `json:"index"`    // "0" (in case of button)
	Parameters []struct {
		Type    string `json:"type"`    // "header", "text", "payload"
		Payload string `json:"payload"` // "hello_world" (in case of payload)
		Text    string `json:"text"`    // "Hello World" (in case of text)
	} `json:"parameters"`
}

func handleSendTemplateMessage(c *fiber.Ctx, template TemplateOptions, to *phonenumber.ParsedPhoneNumber) error {
	if template.Language.Code == "" {
		return customError(c, "(#100) The parameter template['language']['code'] is required.")
	}
	switch strings.ToLower(template.Language.Policy) {
	case "", "deterministic":
		// In case the policy is not set or the value is uppercased
		template.Language.Policy = "deterministic"
	default:
		return customError(c, "(#100) The parameter template['language']['policy'] must be one of {DETERMINISTIC}.")
	}

	msgTemplate := models.Template{}
	err := DB.Model(&models.Template{}).Where("name = ?", template.Name).Preload("TemplateCustomButtons").First(&msgTemplate).Error
	if err != nil {
		msg := "(#132001) Template name does not exist in the translation"
		details := fmt.Sprintf("template name (%s) does not exist in %s", template.Name, template.Language.Code)
		return customError(c, msg, details)
	}

	var requestBodyVariables []string
	var requestHeaderVariables []string
	var buttons []TemplateComponent
	for idx, component := range template.Components {
		switch component.Type {
		case "button":
			buttons = append(buttons, component)
		case "body":
			if requestBodyVariables != nil {
				return customError(c, "There can be at max 1 body component")
			}
			for j, parameter := range component.Parameters {
				if strings.ToLower(parameter.Type) == "text" {
					requestBodyVariables = append(requestBodyVariables, parameter.Text)
				} else {
					msg := fmt.Sprintf("Param template['components'][%d]['parameters'][%d]['type'] must be one of {TEXT}", idx, j)
					return customError(c, msg)
				}
			}
		case "header":
			if requestHeaderVariables != nil {
				return customError(c, "There can be at max 1 header component")
			}
			for j, parameter := range component.Parameters {
				if strings.ToLower(parameter.Type) == "text" {
					requestHeaderVariables = append(requestHeaderVariables, parameter.Text)
				} else {
					msg := fmt.Sprintf("Param template['components'][%d]['parameters'][%d]['type'] must be one of {TEXT}", idx, j)
					return customError(c, msg)
				}
			}
		}
	}

	header := msgTemplate.Header
	body := msgTemplate.Body
	footer := msgTemplate.Footer

	templateBodyVariables := models.Variables(body)
	if len(templateBodyVariables) > 0 {
		if len(requestBodyVariables) != len(templateBodyVariables) {
			msg := "(#132000) Number of parameters does not match the expected number of params"
			detials := fmt.Sprintf(
				"body: number of localizable_params (%d) does not match the expected number of params (%d)",
				len(requestBodyVariables),
				len(templateBodyVariables),
			)
			return customError(c, msg, detials)
		}

		body = models.ReplaceVariables(body, requestBodyVariables)
	}

	if msgTemplate.Header != nil {
		templateHeaderVariables := models.Variables(*msgTemplate.Header)
		if len(templateHeaderVariables) > 0 {
			if len(requestHeaderVariables) != len(templateHeaderVariables) {
				msg := "(#132000) Number of parameters does not match the expected number of params"
				detials := fmt.Sprintf(
					"header: number of localizable_params (%d) does not match the expected number of params (%d)",
					len(requestHeaderVariables),
					len(templateHeaderVariables),
				)
				return customError(c, msg, detials)
			}

			newHeader := models.ReplaceVariables(*header, requestHeaderVariables)
			header = &newHeader
		}
	}

	if len(msgTemplate.TemplateCustomButtons) != len(buttons) {
		msg := "(#132000) Number of parameters does not match the expected number of params"
		details := fmt.Sprintf(
			"number of buttons (%d) does not match the expected number of params (%d)",
			len(buttons),
			len(msgTemplate.TemplateCustomButtons),
		)
		return customError(c, msg, details)
	}

	messageButtons := []models.MessageButton{}
	if buttons != nil {
		type ButtonPayload struct {
			Seen    bool
			Payload string
		}

		buttonsPayload := make([]ButtonPayload, len(buttons))

		for idx, button := range buttons {
			prefix := fmt.Sprintf("template['components'][%d]", idx)

			if button.Index == "" {
				return customError(c, fmt.Sprintf("Param %s['index'] is required", prefix))
			}
			if button.SubType == "" {
				return customError(c, fmt.Sprintf("Param %s['sub_type'] is required", prefix))
			}
			if button.SubType != "quick_reply" {
				return customError(c, fmt.Sprintf("Param %s['sub_type'] must be one of {QUICK_REPLY}", prefix))
			}

			switch len(button.Parameters) {
			case 0:
				return customError(c, fmt.Sprintf("Param %s['parameters'] is required", prefix))
			case 1:
				// continue
			default:
				return customError(c, fmt.Sprintf("Param %s['parameters'] must have at max 1 element", prefix))
			}
			firstParam := button.Parameters[0]
			if firstParam.Type == "" {
				return customError(c, fmt.Sprintf("Param %s['parameters'][0]['type'] is required", prefix))
			}
			if firstParam.Type != "payload" {
				return customError(c, fmt.Sprintf("Param %s['parameters'][0]['type'] must be one of {PAYLOAD}", prefix))
			}
			if firstParam.Payload == "" {
				return customError(c, fmt.Sprintf("Param %s['parameters'][0]['payload'] is required", prefix))
			}

			buttonIndex, err := strconv.Atoi(button.Index)
			if err != nil {
				return customError(c, fmt.Sprintf("Param %s['index'] must be a number", prefix))
			}
			if buttonIndex < 0 || buttonIndex >= len(buttonsPayload) {
				return customError(c, fmt.Sprintf("Param %s['index'] must be between 0 and %d", prefix, len(buttonsPayload)-1))
			}

			buttonsPayload[buttonIndex] = ButtonPayload{
				Seen:    true,
				Payload: firstParam.Payload,
			}

		}

		for idx, btn := range buttonsPayload {
			if !btn.Seen {
				return customError(c, fmt.Sprintf("Button with index %d missing", idx))
			}
			messageButtons = append(messageButtons, models.MessageButton{
				Text:    msgTemplate.TemplateCustomButtons[idx].Text,
				Payload: &btn.Payload,
			})
		}

	}

	message := &models.Message{
		WhatsappID:    to.WhatsappMessageID,
		Direction:     models.DirectionIn,
		HeaderMessage: header,
		Message:       body,
		FooterMessage: footer,
		Timestamp:     time.Now().Unix(),
		Buttons:       messageButtons,
	}
	// Note that templates can be send to everyone
	err = message.CreateOrAppend(to.Parsed)
	if err != nil {
		return err
	}

	websocket.SendMessage(*message)

	// FIXME notify webhook

	return c.JSON(map[string]any{
		"messaging_product": "whatsapp",
		"contacts": []map[string]string{{
			"input": to.Original,
			"wa_id": to.Parsed,
		}},
		"messages": []map[string]string{{
			"id":             message.WhatsappID,
			"message_status": "accepted",
		}},
	})
}
