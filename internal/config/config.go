package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	UseWebhook    bool
	WebhookURL    string
	WebhookSecret string
	Port          int
	Env           string
	LogLevel      string

	PolzaAPIKey string
}

func Load() *Config {

	_ = godotenv.Load(".env")
	_ = godotenv.Load("../../.env")

	portStr := os.Getenv("PORT")
	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 8080
	}

	useWH, _ := strconv.ParseBool(os.Getenv("USE_WEBHOOK"))

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		UseWebhook:    useWH,
		WebhookURL:    os.Getenv("BOT_WEBHOOK_URL"),
		WebhookSecret: os.Getenv("BOT_WEBHOOK_SECRET"),
		Port:          port,
		Env:           os.Getenv("ENV"),
		LogLevel:      firstNonEmpty(os.Getenv("LOG_LEVEL"), "debug"),

		PolzaAPIKey: os.Getenv("POLZA_API_KEY"),
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
