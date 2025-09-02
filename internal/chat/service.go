package chat

import (
	"myapp/pkg/logger"
	"time"

	"go.uber.org/zap"
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
	logger.Log.Info("Sending message",
		zap.String("sessionID", sessionID),
		zap.String("message", message))

	msg := ChatMessage{
		Message:   message,
		SessionID: sessionID,
		Kind:      UserPrompt,
		Timestamp: time.Now().Unix(),
	}
	err := s.repo.Save(&msg)
	if err != nil {
		logger.Log.Error("user message failed to saved", zap.Error(err))
		return Chat{}, err
	}
	messages, err := s.repo.Find(sessionID)
	if err != nil {
		logger.Log.Error("load to history failed", zap.Error(err))
		return Chat{}, err
	}

	response, err := s.client.GetCompletion(message, messages)
	// param := s.client.BuildParams(msg.Message)
	// for _, msg := range messages {
	// 	param.Messages = append(param.Messages, openai.UserMessage(msg.Message))
	// }
	// completion, err = s.client.CreateCompletion(param)

	if err != nil {
		logger.Log.Error("get completion fail", zap.Error(err))
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
		logger.Log.Error("llm response failed to save", zap.Error(err))
		return Chat{}, err
	}

	logger.Log.Info("message sended")
	return Chat{
		Message:   openaiMsg.Message,
		SessionID: openaiMsg.SessionID,
	}, nil
}

func (s *service) FindHistory(sessionID string) ([]ChatMessage, error) {
	logger.Log.Info("Finding history",
		zap.String("sessionID", sessionID))
	messages, err := s.repo.Find(sessionID)
	if err != nil {
		logger.Log.Error("failed to load history", zap.Error(err))
		return nil, err
	}
	logger.Log.Info("history loaded")
	return messages, nil
}
