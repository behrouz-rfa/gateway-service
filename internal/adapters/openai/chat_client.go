package chat

import (
	"context"
	"errors"
	openai "github.com/sashabaranov/go-openai"
)

const DefaultAi = openai.GPT4o

// Define an interface for the OpenAI client
type OpenAIClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error)
}

// Adapter struct that wraps the openai.Client
type OpenAIClientAdapter struct {
	client *openai.Client
}

func (a *OpenAIClientAdapter) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	resp, err := a.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update Client struct to use the interface
type Client struct {
	AiClient OpenAIClient
}

func NewClient(key string) *Client {
	client := new(Client)
	client.AiClient = &OpenAIClientAdapter{client: openai.NewClient(key)}

	return client
}

func (c *Client) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (string, error) {
	if c.AiClient == nil {
		return "", nil
	}

	req.Model = DefaultAi
	resp, err := c.AiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("empty")
	}
	content := ""
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Message.Content
	}

	return content, nil
}
