package main

import (
	"context"
	"log"

	"github.com/go-telegram/bot"

	"ai-telegram-saas/internal/application/service"
	"ai-telegram-saas/internal/config"
	polza "ai-telegram-saas/internal/infrastructure/ai/provider/polza"
	"ai-telegram-saas/internal/infrastructure/telegram/handlers" // ПРОВЕРЬ ЭТОТ ПУТЬ
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	aiProvider := polza.New(cfg.PolzaAPIKey)
	chatService := service.NewChatService(aiProvider)

	// Передаем зависимости правильно
	deps := &handlers.Deps{
		Chat: chatService,
	}

	b, err := bot.New(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypeContains,
		handlers.Echo(deps),
	)

	log.Println("BOT STARTED")
	b.Start(ctx)
}
