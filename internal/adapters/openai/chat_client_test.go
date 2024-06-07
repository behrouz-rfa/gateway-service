package chat

import (
	"context"
	"errors"
	"testing"

	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

// Mock client implementing the OpenAIClient interface
type MockOpenAIClient struct {
	mockCreateChatCompletion func(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error)
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	return m.mockCreateChatCompletion(ctx, req)
}

func TestCreateChatCompletion(t *testing.T) {
	tests := []struct {
		name        string
		req         openai.ChatCompletionRequest
		mockResp    *openai.ChatCompletionResponse
		mockErr     error
		expected    string
		expectError bool
	}{
		{
			name: "Successful chat completion",
			req:  openai.ChatCompletionRequest{},
			mockResp: &openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{Message: openai.ChatCompletionMessage{Content: "Hello, world!"}},
				},
			},
			expected:    "Hello, world!",
			expectError: false,
		},
		{
			name:        "Client is nil",
			req:         openai.ChatCompletionRequest{},
			mockResp:    nil,
			mockErr:     nil,
			expected:    "",
			expectError: true,
		},
		{
			name:        "API call fails",
			req:         openai.ChatCompletionRequest{},
			mockResp:    nil,
			mockErr:     errors.New("API call failed"),
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockOpenAIClient{
				mockCreateChatCompletion: func(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
					return tt.mockResp, tt.mockErr
				},
			}

			client := &Client{AiClient: mockClient}

			resp, err := client.CreateChatCompletion(context.Background(), tt.req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp)
			}
		})
	}
}
