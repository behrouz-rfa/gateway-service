package ports

import "github.com/behrouz-rfa/gateway-service/internal/core/domain"

type PlanRepository interface {
	GetUserPlan(userID string) (*domain.Plan, error)
	UpdateUserCredits(userID string, creditsUsed int) error
}
