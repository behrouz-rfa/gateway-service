package db

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
)

func (r *DbRepository) CreateUserPlan(ctx context.Context, plan *domain.Plan) (string, error) {

	if err := r.db.WithContext(ctx).Create(plan).Error; err != nil {
		return "", err
	}

	return plan.ID, nil
}

func (r *DbRepository) GetUserPlan(userID string) (*domain.Plan, error) {
	var userPlan domain.Plan
	if err := r.db.Model(&domain.Plan{}).Where("user_id = ?", userID).First(&userPlan).Error; err != nil {
		return nil, err
	}
	return &userPlan, nil
}

func (r *DbRepository) UpdateUserCredits(userID string, creditsUsed int) error {
	return r.db.Model(&domain.Plan{}).Where("user_id = ?", userID).Update("credits_used", creditsUsed).Error
}

func (r *DbRepository) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *DbRepository) CreateUser(ctx context.Context, user *domain.User) (string, error) {

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return "", err
	}

	return user.ID, nil
}

func (r *DbRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var usr domain.User
	if err := r.db.Where("email = ?", email).First(&usr).Error; err != nil {
		return nil, err
	}
	return &usr, nil
}
