package chat

import (
	"myapp/internal/llm"
	"myapp/internal/models"
	"time"
)

type Service interface {
	SendMessage(sessionID string, message string) (models.Chat, error)
	FindHistory(sessionID string) ([]models.ChatMessage, error)
}

type service struct {
	repo   Repository
	client llm.Client
}

func NewService(repo Repository, llmClient llm.Client) Service {
	return &service{
		repo:   repo,
		client: llmClient,
	}
}
func (s *service) SendMessage(sessionID string, message string) (models.Chat, error) {
	msg := models.ChatMessage{
		Message:   message,
		SessionID: sessionID,
		Kind:      models.UserPrompt,
		Timestamp: time.Now().Unix(),
	}
	err := s.repo.Save(&msg)
	if err != nil {
		return models.Chat{}, err
	}
	messages, err := s.repo.Find(sessionID)
	if err != nil {
		return models.Chat{}, err
	}

	response, err := s.client.GetCompletion(message, messages)
	// param := s.client.BuildParams(msg.Message)
	// for _, msg := range messages {
	// 	param.Messages = append(param.Messages, openai.UserMessage(msg.Message))
	// }
	// completion, err = s.client.CreateCompletion(param)

	if err != nil {
		return models.Chat{}, err
	}
	openaiMsg := models.ChatMessage{
		Message:   response,
		SessionID: sessionID,
		Kind:      models.LLMOutput,
		Timestamp: time.Now().Unix(),
	}
	err = s.repo.Save(&openaiMsg)
	if err != nil {
		return models.Chat{}, err
	}

	return models.Chat{
		Message:   openaiMsg.Message,
		SessionID: openaiMsg.SessionID,
	}, nil
}

func (s *service) FindHistory(sessionID string) ([]models.ChatMessage, error) {
	messages, err := s.repo.Find(sessionID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
