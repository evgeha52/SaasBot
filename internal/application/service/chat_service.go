package service

import (
	"context"
	"fmt"
	"sync"
)

type AIProvider interface {
	Ask(ctx context.Context, model string, message string) (string, error)
}

type ChatService struct {
	ai         AIProvider
	models     []string
	modelMap   map[int64]string
	historyMap map[int64][]string // Храним историю просто как строки
	mu         sync.RWMutex
}

func NewChatService(provider AIProvider) *ChatService {
	return &ChatService{
		ai:         provider,
		models:     []string{"deepseek-chat", "gpt-4o-mini", "gpt-4o"},
		modelMap:   make(map[int64]string),
		historyMap: make(map[int64][]string),
	}
}

func (s *ChatService) Models() []string { return s.models }

func (s *ChatService) SetModel(chatID int64, model string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.modelMap[chatID] = model
}

func (s *ChatService) Ask(ctx context.Context, chatID int64, text string) (string, error) {
	s.mu.Lock()
	model, ok := s.modelMap[chatID]
	if !ok {
		model = "deepseek-chat"
	}

	// Склеиваем историю в один большой текст
	fullPrompt := ""
	for _, msg := range s.historyMap[chatID] {
		fullPrompt += msg + "\n"
	}
	fullPrompt += "User: " + text
	s.mu.Unlock()

	// Отправляем просто текст провайдеру
	answer, err := s.ai.Ask(ctx, model, fullPrompt)
	if err != nil {
		return "", err
	}

	// Сохраняем в историю
	s.mu.Lock()
	s.historyMap[chatID] = append(s.historyMap[chatID], fmt.Sprintf("User: %s", text))
	s.historyMap[chatID] = append(s.historyMap[chatID], fmt.Sprintf("Assistant: %s", answer))

	// Держим только последние 10 сообщений
	if len(s.historyMap[chatID]) > 10 {
		s.historyMap[chatID] = s.historyMap[chatID][len(s.historyMap[chatID])-10:]
	}
	s.mu.Unlock()

	return answer, nil
}

func (s *ChatService) Reset(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.historyMap, chatID)
}
