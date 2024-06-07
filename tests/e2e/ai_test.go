//go:build e2e

package gql

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/common"
	"github.com/goccy/go-json"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *GqlTestSuite) TestAiUser() {
	// Register User
	inputReq := `{
		"name": "Jon",
		"password": "12345678",
		"email": "joa@gmail.com"
	}`
	req, err := http.NewRequest("POST", "/api/v1/users/register", strings.NewReader(inputReq))
	if err != nil {
		s.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.api.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var userData userResponseData
	err = json.Unmarshal(w.Body.Bytes(), &userData)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(userData.Data.Email, "joa@gmail.com")

	// Interact with OpenAI API
	aiRequest := `{
		"content": "Jon"
	}`
	req2, err := http.NewRequest("POST", "/api/v1/openai", strings.NewReader(aiRequest))
	if err != nil {
		s.T().Fatal(err)
	}
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set(string(common.AuthorizationContextKey), "Bearer "+userData.Data.JWTToken.Token)

	w = httptest.NewRecorder()
	s.api.ServeHTTP(w, req2)
	s.Equal(http.StatusOK, w.Code)

	var aiData aiResponse
	err = json.Unmarshal(w.Body.Bytes(), &aiData)
	if err != nil {
		s.T().Fatal(err)
	}

}

// Mock client implementing the OpenAIClient interface
type MockOpenAIClient struct {
	MockCreateChatCompletion func(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error)
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	return m.MockCreateChatCompletion(ctx, req)
}
