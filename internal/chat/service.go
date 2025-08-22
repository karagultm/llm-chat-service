package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/openai/openai-go/v2"
)

type Service interface {
	SendMessage(sessionID string, message string) (Chat, error)
	FindHistory(sessionID string) ([]ChatMessage, error)
}

type service struct {
	repo   Repository
	client *openai.Client
}

func NewService(repo Repository, client *openai.Client) Service {
	return &service{
		repo:   repo,
		client: client,
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

	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(msg.Message),
		},
		Seed:  openai.Int(1),
		Model: openai.ChatModelGPT4o,
	}
	messages, _ := s.repo.Find(sessionID) //error  eklenirs
	for _, msg := range messages {
		param.Messages = append(param.Messages, openai.UserMessage(msg.Message))
		fmt.Println(msg.Message)
	}
	completion, err := s.client.Chat.Completions.New(context.TODO(), param)

	if err != nil {
		return Chat{}, err
	}
	openaiMsg := ChatMessage{
		Message:   completion.Choices[0].Message.Content,
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
		return []ChatMessage{}, err
	}
	return messages, nil
}
