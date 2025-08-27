package llm

import (
	"context"
	"fmt"
	"myapp/internal/models"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type Client interface {
	GetCompletion(message string, messages []models.ChatMessage) (response string, err error)
}

type client struct {
	openai openai.Client
}

func NewClient(apiKey string) Client {
	return &client{
		openai: (openai.NewClient(option.WithAPIKey(apiKey))),
	}
}

func (c *client) GetCompletion(message string, messages []models.ChatMessage) (response string, err error) {
	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
		Seed:  openai.Int(1),
		Model: openai.ChatModelGPT4o,
	}
	for _, msg := range messages {
		switch msg.Kind {
		case models.UserPrompt:
			param.Messages = append(param.Messages, openai.UserMessage(msg.Message))
		case models.LLMOutput:
			param.Messages = append(param.Messages, openai.AssistantMessage(msg.Message))
		}
	}

	completion, err := c.openai.Chat.Completions.New(context.TODO(), param)
	if err != nil {
		return "", err
	}

	// completion doluysa response'u al
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no choices returned by OpenAI")
	}

	response = completion.Choices[0].Message.Content
	return
}
