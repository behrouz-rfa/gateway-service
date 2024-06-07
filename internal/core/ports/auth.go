package ports

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
)

type AuthRepository interface {
	Create(info domain.UserClaims) (*domain.JWTToken, error)
	Verify(token string) (*domain.UserClaims, error)
	Authorize(ctx context.Context, token string) (bool, error)
	GetUserClaims(ctx context.Context, token string) (*domain.UserClaims, error)
	GetUserClaimsByEmail(ctx context.Context, email string) (*domain.UserClaims, error)
}

type AuthService interface {
	Create(info domain.UserClaims) (*domain.JWTToken, error)
	Verify(token string) (*domain.UserClaims, error)
	Authorize(ctx context.Context, token string) (bool, error)
	GetUserClaims(ctx context.Context, token string) (*domain.UserClaims, error)
	GetUserClaimsByEmail(ctx context.Context, email string) (*domain.UserClaims, error)
	GetUser(ctx context.Context, token string) (*domain.User, error)
}
