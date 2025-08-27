package chat

import (
	"time"
)

type Service interface {
	SendMessage(sessionID string, message string) (Chat, error)
	FindHistory(sessionID string) ([]ChatMessage, error)
}

type service struct {
	repo   Repository
	client Client
}

func NewService(repo Repository, llmClient Client) Service {
	return &service{
		repo:   repo,
		client: llmClient,
	}
}
func (s *service) SendMessage(sessionID string, message string) (Chat, error) {
	msg := ChatMessage{
		Message:   message,
		SessionID: sessionID,
		Kind:      UserPrompt,
		Timestamp: time.Now().Unix(),
	}
	err := s.repo.Save(&msg)
	if err != nil {
		return Chat{}, err
	}
	messages, err := s.repo.Find(sessionID)
	if err != nil {
		return Chat{}, err
	}

	response, err := s.client.GetCompletion(message, messages)
	// param := s.client.BuildParams(msg.Message)
	// for _, msg := range messages {
	// 	param.Messages = append(param.Messages, openai.UserMessage(msg.Message))
	// }
	// completion, err = s.client.CreateCompletion(param)

	if err != nil {
		return Chat{}, err
	}
	openaiMsg := ChatMessage{
		Message:   response,
		SessionID: sessionID,
		Kind:      LLMOutput,
		Timestamp: time.Now().Unix(),
	}
	err = s.repo.Save(&openaiMsg)
	if err != nil {
		return Chat{}, err
	}

	return Chat{
		Message:   openaiMsg.Message,
		SessionID: openaiMsg.SessionID,
	}, nil
}

func (s *service) FindHistory(sessionID string) ([]ChatMessage, error) {
	messages, err := s.repo.Find(sessionID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
