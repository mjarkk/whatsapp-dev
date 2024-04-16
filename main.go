package main

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"

	. "github.com/mjarkk/whatsapp-dev/go"
	. "github.com/mjarkk/whatsapp-dev/go/db"
	"github.com/mjarkk/whatsapp-dev/go/lib/webhook"
	"github.com/mjarkk/whatsapp-dev/go/models"
	"github.com/mjarkk/whatsapp-dev/go/state"
	"github.com/mjarkk/whatsapp-dev/go/utils/random"

	"github.com/spf13/pflag"
)

//go:embed dist
var dist embed.FS

func getenv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

type ArgOrEnvPromise func() string

func argOrEnv(flagLong, flagShort, envKey, defaultValue, description string) ArgOrEnvPromise {
	arg := pflag.StringP(flagLong, flagShort, defaultValue, description)
	return func() string {
		env := getenv(envKey)
		if *arg == defaultValue && env != "" {
			*arg = env
		}
		return *arg
	}
}

func main() {
	defaultHttpAddr := ":1090"

	webHookURL := argOrEnv("webhook-url", "w", "WEBHOOK_URL", "", "Webhook URL")
	webHookVerivyToken := argOrEnv("webhook-verify-token", "t", "WEBHOOK", "", "Webhook verify token")
	secretesSeed := argOrEnv("secrets-seed", "", "SECRETS_SEED", "", "Secrets seed for generating a random graph-token")
	httpAddr := argOrEnv("http-addr", "a", "HTTP_ADDR", defaultHttpAddr, "HTTP address")
	httpUsername := argOrEnv("http-username", "u", "HTTP_USERNAME", "", "HTTP username")
	httpPassword := argOrEnv("http-password", "p", "HTTP_PASSWORD", "", "HTTP password")
	phoneNumber := argOrEnv("whatsapp-phone-number", "", "WHATSAPP_PHONE_NUMBER", "", "Define the mocked phone number")
	phoneNumberID := argOrEnv("whatsapp-phone-number-id", "", "WHATSAPP_PHONE_NUMBER_ID", "", "Define the mocked phone number id")
	graphToken := argOrEnv("facebook-graph-token", "", "FACEBOOK_GRAPH_TOKEN", "", "Define mock graph token")
	appSecret := argOrEnv("facebook-app-secret", "", "FACEBOOK_APP_SECRET", "", "Define the mocked phone number id")

	pflag.Parse()

	webHookURLValue := webHookURL()
	if webHookURLValue == "" {
		panic("Webhook url must be set using the --webhook-url flag or the $WEBHOOK_URL environment variable")
	}
	_, err := url.Parse(webHookURLValue)
	if err != nil {
		panic("Invalid webhook url: " + err.Error())
	}

	secretesSeedValue := secretesSeed()
	if secretesSeedValue == "" {
		fmt.Println("DANGER: using fallback secrets seed")
		fmt.Println("        not recommended when exposing this service to the internet")
		fmt.Println("        use --secrets-seed or $SECRETS_SEED to set a custom seed")
		secretesSeedValue = "fallback-secrets-seed"
	}

	h := sha256.New()
	h.Write([]byte(secretesSeedValue))
	result := h.Sum(nil)
	var seed int64
	for idx, b := range result {
		seed ^= int64(b) << idx * 2
	}

	r := rand.New(rand.NewSource(seed))
	initialRandomValues := random.GetRandomValuesForSetup(r)

	graphTokenValue := graphToken()
	if graphTokenValue == "" {
		graphTokenValue = initialRandomValues.GraphToken
	}

	appSecretValue := appSecret()
	if appSecretValue == "" {
		appSecretValue = initialRandomValues.AppSecret
	}

	phoneNumberValue := phoneNumber()
	if phoneNumberValue == "" {
		phoneNumberValue = initialRandomValues.PhoneNumber
	}

	phoneNumberIDValue := phoneNumberID()
	if phoneNumberIDValue == "" {
		phoneNumberIDValue = initialRandomValues.PhoneNumberID
	}

	webhookVerifyTokenValue := webHookVerivyToken()
	if webhookVerifyTokenValue == "" {
		webhookVerifyTokenValue = initialRandomValues.WebhookVerifyToken
	}

	fmt.Println("Graph token:\t", graphTokenValue)
	state.GraphToken.Set(graphTokenValue)
	fmt.Println("App secret:\t", appSecretValue)
	state.AppSecret.Set(appSecretValue)
	fmt.Println("Phone number:\t", phoneNumberValue)
	state.PhoneNumber.Set(phoneNumberValue)
	fmt.Println("Phone number ID:", phoneNumberIDValue)
	state.PhoneNumberID.Set(phoneNumberIDValue)
	fmt.Println("Webhook verify token:", webhookVerifyTokenValue)
	state.WebhookVerifyToken.Set(webhookVerifyTokenValue)

	state.WebhookURL.Set(webHookURLValue)

	ConnectToDatabase()

	DB.AutoMigrate(
		&models.Conversation{},
		&models.Message{},
		&models.Template{},
		&models.TemplateCustomButton{},
	)

	templatesCount := int64(0)
	err = DB.Model(&models.Template{}).Count(&templatesCount).Error
	if err != nil {
		panic(err)
	}

	if templatesCount == 0 {
		header := "Hello World"
		footer := "WhatsApp dev sample message"
		DB.Create(&models.Template{
			Name:   "hello_world",
			Body:   "Welcome and congratulations!! This message demonstrates your ability to send a WhatsApp message notification from the Cloud API, hosted by whatsapp dev. Thank you for taking the time to test with us.",
			Header: &header,
			Footer: &footer,
		})
	}

	go func() {
		err := webhook.Validate()
		if err != nil {
			fmt.Println("Failed to validate webhook:", err.Error())
		}
	}()

	StartWebserver(StartWebserverOptions{
		Addr:              httpAddr(),
		BasicAuthUsername: httpUsername(),
		BasicAuthPassword: httpPassword(),
		Rand:              r,
		Dist:              dist,
	})
}
