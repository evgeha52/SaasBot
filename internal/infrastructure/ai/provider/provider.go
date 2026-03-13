package provider

import "context"

type Message struct {
	Role    string
	Content string
}

type AIProvider interface {
	ChatWithHistory(ctx context.Context, model string, history []Message) (string, error)
}
