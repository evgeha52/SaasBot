package polza

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type Provider struct {
	client *openai.Client
}

func New(apiKey string) *Provider {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = "https://api.polza.ai/v1"
	client := openai.NewClientWithConfig(cfg)

	return &Provider{client: client}
}

func (p *Provider) Ask(ctx context.Context, model string, message string) (string, error) {
	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "STRICT RULES: Use ONLY Telegram HTML. " +
						"Wrap ALL code and formulas in <pre>...</pre>. " +
						"Use <b>...</b> for bold headers. " +
						"NO MARKDOWN (#, **, ```).",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
