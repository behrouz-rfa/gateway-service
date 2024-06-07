package ports

import "context"

type OpenAIService interface {
	GenerateText(ctx context.Context, prompt string, userID string) (string, error)
}
