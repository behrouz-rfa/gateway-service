package ports

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, filmInput *domain.User) (string, error)
	CreateUserPlan(ctx context.Context, filmInput *domain.Plan) (string, error)
}

type UserService interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, filmInput *domain.UserInput) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	Login(ctx context.Context, input *domain.UserLoginInput) (*domain.User, error)
}
