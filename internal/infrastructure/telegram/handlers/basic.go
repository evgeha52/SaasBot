package handlers

import (
	"context"
	"strings"

	"ai-telegram-saas/internal/application/service"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Deps struct {
	Chat *service.ChatService
}

// форматирование ответа AI
func formatOutput(text string) string {

	lines := strings.Split(text, "\n")

	var result strings.Builder
	inCode := false

	for _, line := range lines {

		// начало / конец блока кода
		if strings.HasPrefix(line, "```") {

			if !inCode {
				result.WriteString("<pre><code>")
				inCode = true
			} else {
				result.WriteString("</code></pre>\n")
				inCode = false
			}

			continue
		}

		if inCode {

			// экранируем html
			line = strings.ReplaceAll(line, "&", "&amp;")
			line = strings.ReplaceAll(line, "<", "&lt;")
			line = strings.ReplaceAll(line, ">", "&gt;")

			result.WriteString(line + "\n")

		} else {

			// убираем markdown мусор
			line = strings.ReplaceAll(line, "**", "")
			line = strings.ReplaceAll(line, "*", "")
			line = strings.ReplaceAll(line, "###", "")
			line = strings.ReplaceAll(line, "##", "")
			line = strings.ReplaceAll(line, "#", "")

			result.WriteString(line + "\n")
		}
	}

	if inCode {
		result.WriteString("</code></pre>")
	}

	return result.String()
}

func Echo(deps *Deps) bot.HandlerFunc {

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

		if update.Message == nil {
			return
		}

		text := strings.TrimSpace(update.Message.Text)
		chatID := update.Message.Chat.ID

		if text == "" {
			return
		}

		// список моделей
		if text == "/models" {

			list := deps.Chat.Models()

			msg := "<b>Доступные модели:</b>\n\n"

			for _, m := range list {
				msg += "• <code>" + m + "</code>\n"
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      msg,
				ParseMode: models.ParseModeHTML,
			})

			return
		}

		// запрос к AI
		answer, err := deps.Chat.Ask(ctx, chatID, text)

		if err != nil {
			answer = "Ошибка: " + err.Error()
		}

		finalText := formatOutput(answer)

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      finalText,
			ParseMode: models.ParseModeHTML,
		})

		if err != nil {

			// fallback если html сломался
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   answer,
			})

		}
	}
}
