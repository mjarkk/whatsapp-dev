package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	. "github.com/mjarkk/whatsapp-dev/go/db"
	"github.com/mjarkk/whatsapp-dev/go/models"
	"github.com/mjarkk/whatsapp-dev/go/state"
	"github.com/mjarkk/whatsapp-dev/go/utils/random"
)

var httpClient = http.Client{
	Timeout: time.Second * 30,
}

type M map[string]any

func createSignatures(body []byte) map[string]string {
	appSecret := state.AppSecret.Get()

	sha1Hmac := hmac.New(sha1.New, []byte(appSecret))
	sha1Hmac.Write(body)
	sha1Signature := "sha1=" + string(sha1Hmac.Sum(nil))

	sha256Hmac := hmac.New(sha256.New, []byte(appSecret))
	sha256Hmac.Write(body)
	sha256Signature := "sha256=" + string(sha256Hmac.Sum(nil))

	return map[string]string{
		"x-hub-signature":     sha1Signature,
		"x-hub-signature-256": sha256Signature,
	}
}

func Validate() error {
	url, err := url.Parse(state.WebhookURL.Get())
	if err != nil {
		return err
	}

	randomSource := rand.New(rand.NewSource(time.Now().Unix()))
	challenge := random.Hex(randomSource, 16)

	url.Query().Add("hub.mode", "subscribe")
	url.Query().Add("hub.verivy_token", state.WebhookVerifyToken.Get())
	url.Query().Add("hub.challenge", challenge)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respBodyStr := string(respBody)
	if respBodyStr == challenge {
		// Yay success
		return nil
	}

	return fmt.Errorf(
		"webhook did not return the expected challenge, expected: \"%s\", got: \"%s\"",
		challenge,
		respBodyStr,
	)
}

func NotivyMessage(message models.Message, awaitResponse bool) error {
	conversation := models.Conversation{}
	err := DB.Model(models.Conversation{}).Find(&conversation, message.ConversationID).Error
	if err != nil {
		return err
	}

	data := M{
		"object": "whatsapp_business_account",
		"entry": []M{{
			"id": strconv.Itoa(int(message.ID)),
			"changes": []M{{
				"value": M{
					"messaging_product": "whatsapp",
					"metadata": M{
						"display_phone_number": state.PhoneNumber.Get(),
						"phone_number_id":      state.PhoneNumberID.Get(),
					},
					"contacts": []M{{
						// FIXME Add a custom contact name to the conversation
						// FIXME I think the production whatsapp api version does not always send the contact name
						"profile": M{"name": "Jhon doe"},
						"wa_id":   conversation.PhoneNumberId,
					}},
					"messages": []M{{
						"from":      conversation.PhoneNumber,
						"id":        message.WhatsappID,
						"timestamp": strconv.Itoa(int(message.Timestamp)),
						"text":      M{"body": message.Message},
						"type":      "text",
					}},
				},
				"field": "messages",
			}},
		}},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if !awaitResponse {
		go makeWebhookRequest(payload)
		return nil
	}

	return makeWebhookRequest(payload)
}

func makeWebhookRequest(payload []byte) error {
	attempt := 0
	var lastErr error
	for {
		attempt++
		switch attempt {
		case 1:
			// Just continue
		case 2:
			time.Sleep(time.Second * 2)
		case 3:
			time.Sleep(time.Second * 5)
		case 4:
			time.Sleep(time.Second * 15)
		default:
			break
		}

		body := bytes.NewBuffer(payload)
		req, err := http.NewRequest("POST", state.WebhookURL.Get(), body)
		if err != nil {
			lastErr = err
			continue
		}

		signatures := createSignatures(payload)
		for key, value := range signatures {
			req.Header.Add(key, value)
		}
		req.Header.Add("user-agent", "facebookexternalua")
		req.Header.Add("content-type", "application/json")

		if attempt == 1 {
			fmt.Println("calling webhook")
		} else {
			fmt.Println("retrying webhook")
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		resp.Body.Close()

		if resp.StatusCode < 400 {
			return nil
		}
	}

	return lastErr
}
