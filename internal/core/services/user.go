package services

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"
	util "github.com/behrouz-rfa/gateway-service/pkg/utils"
)

// UserService is a service that handles user-related operations.
type UserServiceOpts struct {
	UserRepo   ports.UserRepository
	Auth       ports.AuthService
	PlanCredit int
}

// UserService is a service that handles user-related operations.
type UserService struct {
	UserServiceOpts

	lg *logger.Entry
}

// UserServiceOption is a function that configures the UserService.
type UserServiceOption func(*UserService)

// NewUserService creates a new instance of UserService with the provided options.
func NewUserService(opts UserServiceOpts) *UserService {
	s := &UserService{
		UserServiceOpts: opts,
		lg:              logger.General.Component("UserService"),
	}

	return s
}

// GetUser retrieves a user by the specified filter.
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {

	return s.UserRepo.GetUser(ctx, id)
}

// CreateUser creates a new user with the provided input.
func (s *UserService) CreateUser(ctx context.Context, input *domain.UserInput) (*domain.User, error) {

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	userID, err := s.UserRepo.CreateUser(ctx, &domain.User{
		Password: hashedPassword,
		Name:     input.Name,
		Email:    input.Email,
	})

	if err != nil {
		return nil, err
	}

	_, err = s.UserRepo.CreateUserPlan(ctx, &domain.Plan{
		UserID:      userID,
		CreditsUsed: 0,
		CreditLimit: s.PlanCredit,
	})

	if err != nil {
		return nil, err
	}

	token, err := s.Auth.Create(domain.UserClaims{UserID: userID, Email: input.Email})
	if err != nil {
		return nil, domain.ErrInternal
	}

	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.JWTToken = token

	return user, nil
}

// Login authenticates a user with the provided login input.
func (s *UserService) Login(ctx context.Context, input *domain.UserLoginInput) (*domain.User, error) {

	user, err := s.UserRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := util.ComparePassword(input.Password, user.Password); err != nil {
		return nil, domain.ErrInternal
	}

	if err != nil {
		return nil, err
	}

	token, err := s.Auth.Create(domain.UserClaims{Email: user.Email, UserID: user.ID})
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.JWTToken = token
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
