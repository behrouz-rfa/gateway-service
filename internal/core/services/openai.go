package services

import (
	chat "github.com/behrouz-rfa/gateway-service/internal/adapters/openai"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/net/context"
	"net/http"
)

type OpenAIServiceOpts struct {
	Gpt       *chat.Client
	PlanRepo  ports.PlanRepository
	UserRepo  ports.UserRepository
	RedisRepo ports.RedisRepository
}

type OpenAIService struct {
	opts OpenAIServiceOpts
	lg   *logger.Entry
}

func NewOpenAIService(opts OpenAIServiceOpts) *OpenAIService {
	return &OpenAIService{
		opts: opts,
		lg:   logger.General.Component("OpenAiService"),
	}
}

func (ai *OpenAIService) GenerateText(ctx context.Context, param string, userID string) (string, error) {
	userPlan, err := ai.opts.PlanRepo.GetUserPlan(userID)
	if err != nil {
		ai.lg.WithError(err).Info("failed to retrieve user plan")
		return "", domain.ErrDataNotFound
	}
	if userPlan.CreditsUsed >= userPlan.CreditLimit {
		ai.lg.WithError(domain.ErrInsufficientPayment).Info("insufficient credits")
		return "", domain.ErrInsufficientPayment
	}

	if response, ok := ai.opts.RedisRepo.GetCachedResponse(param); ok {
		return string(response), nil
	}

	completion, err := ai.opts.Gpt.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Messages: []openai.ChatCompletionMessage{
			{
				Content: param,
				Role:    openai.ChatMessageRoleUser,
			},
		},
	})
	if err != nil {
		ai.lg.WithError(err).Info("CreateChatCompletion error occurred")
		return "", &domain.HTTPError{
			Message:    "failed to generate text",
			StatusCode: http.StatusInternalServerError,
		}
	}

	go ai.updateUserCreditsAndCacheResponse(userID, param, userPlan.CreditsUsed, completion)

	return completion, nil
}

func (ai *OpenAIService) updateUserCreditsAndCacheResponse(userID, param string, creditsUsed int, completion string) {
	if len(completion) > 0 {
		err := ai.opts.PlanRepo.UpdateUserCredits(userID, creditsUsed+1)
		if err != nil {
			ai.lg.WithError(err).Info("failed to update user credits")
		}
		ai.opts.RedisRepo.SetCachedResponse(param, []byte(completion))
	}
}
