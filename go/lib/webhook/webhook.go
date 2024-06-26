package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
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
	sha1Signature := "sha1=" + hex.EncodeToString(sha1Hmac.Sum(nil))

	sha256Hmac := hmac.New(sha256.New, []byte(appSecret))
	sha256Hmac.Write(body)
	sha256Signature := "sha256=" + hex.EncodeToString(sha256Hmac.Sum(nil))

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

	query := url.Query()
	query.Add("hub.mode", "subscribe")
	query.Add("hub.verify_token", state.WebhookVerifyToken.Get())
	query.Add("hub.challenge", challenge)
	url.RawQuery = query.Encode()

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

	timestamp := strconv.Itoa(int(message.Timestamp))
	from := conversation.PhoneNumber
	bodyMessage := M{
		"from":      from,
		"id":        message.WhatsappID,
		"timestamp": timestamp,
		"text":      M{"body": message.Message},
		"type":      "text",
	}
	if message.Payload != nil {
		bodyMessage = M{
			"from":      from,
			"id":        message.WhatsappID,
			"timestamp": timestamp,
			"type":      "button",
			"button": M{
				"payload": *message.Payload,
				"text":    message.Message,
			},
		}
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
					"messages": []M{bodyMessage},
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
		go func() {
			err := makeWebhookRequestWithRandomForcedRetries(payload)
			if err != nil {
				fmt.Println("failed to call webhook, error response:", err.Error())
			}
		}()
		return nil
	}

	return makeWebhookRequestWithRandomForcedRetries(payload)
}

func makeWebhookRequestWithRandomForcedRetries(payload []byte) error {
	randomSleepDuration := rand.Intn(int(time.Millisecond * 1500))
	time.Sleep(time.Duration(randomSleepDuration))

	err := makeWebhookRequest(payload)
	if err != nil {
		return err
	}

	if rand.Intn(100) > 20 {
		return nil
	}
	// Yolo retry 20% of all request.
	// Facebooks api also does this and can be a source of bugs hence why we also do this.

	randomSleepDuration = rand.Intn(int(time.Second * 10))
	time.Sleep(time.Duration(randomSleepDuration))

	err = makeWebhookRequest(payload)
	if err != nil {
		return err
	}

	if rand.Intn(100) > 10 {
		return nil
	}
	// Yolo retry 10% of all request (this should be 2% of the time in total).
	// Same reason as above

	randomSleepDuration = rand.Intn(int(time.Minute))
	time.Sleep(time.Duration(randomSleepDuration))

	return makeWebhookRequest(payload)
}

func makeWebhookRequest(payload []byte) error {
	attempt := 0
	var lastErr error

outer:
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
			break outer
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
