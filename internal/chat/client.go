package chat

import (
	"context"
	"fmt"
	"myapp/pkg/logger"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"go.uber.org/zap"
)

type Client interface {
	GetCompletion(message string, messages []ChatMessage) (response string, err error)
}

type client struct {
	openai openai.Client
}

func NewClient(apiKey string) Client {
	return &client{
		openai: (openai.NewClient(option.WithAPIKey(apiKey))),
	}
}

func (c *client) GetCompletion(message string, messages []ChatMessage) (response string, err error) {
	logger.Log.Info("Client received user message",
		zap.String("message", message))
	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
		Seed:  openai.Int(1),
		Model: openai.ChatModelGPT4o,
	}
	for _, msg := range messages {
		switch msg.Kind {
		case UserPrompt:
			param.Messages = append(param.Messages, openai.UserMessage(msg.Message))
		case LLMOutput:
			param.Messages = append(param.Messages, openai.AssistantMessage(msg.Message))
		}
	}

	completion, err := c.openai.Chat.Completions.New(context.TODO(), param)
	if err != nil {
		logger.Log.Error("Failed to get OpenAI completion",
			zap.String("message", message),
			zap.Any("chat_history", messages),
			zap.Error(err))
		return "", err
	}

	// completion doluysa response'u al
	if len(completion.Choices) == 0 {
		logger.Log.Warn("OpenAI returned empty choices",
			zap.String("message", message),
			zap.Any("chat_history", messages))
		return "", fmt.Errorf("no choices returned by OpenAI")
	}

	response = completion.Choices[0].Message.Content
	logger.Log.Info("OpenAI responed successfully",
		zap.String("response", response),
		zap.String("message", message))
	return
}
