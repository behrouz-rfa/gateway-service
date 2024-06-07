package services

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"
	"net/http"
)

type GeolocationServiceOpts struct {
	PlanRepo  ports.PlanRepository
	Repo      ports.GeolocationRepository
	RedisRepo ports.RedisRepository
}

type GeolocationService struct {
	opts GeolocationServiceOpts
	lg   *logger.Entry
}

func NewGeolocationService(opts GeolocationServiceOpts) *GeolocationService {
	return &GeolocationService{
		opts: opts,
		lg:   logger.General.Component("GeolocationService"),
	}
}

func (s *GeolocationService) GetGeolocation(ctx context.Context, ip string, userID string) (*domain.GeolocationResponse, error) {
	userPlan, err := s.opts.PlanRepo.GetUserPlan(userID)
	if err != nil {
		s.lg.WithError(err).Info("failed to retrieve user plan")
		return nil, &domain.HTTPError{
			Message:    "failed to retrieve user plan",
			StatusCode: http.StatusInternalServerError,
		}
	}

	if userPlan.CreditsUsed >= userPlan.CreditLimit {
		s.lg.WithError(domain.ErrInsufficientPayment).Info("insufficient credits")
		return nil, domain.ErrInsufficientPayment
	}

	geolocationResponse := new(domain.GeolocationResponse)
	err = s.opts.RedisRepo.Get(ctx, ip, geolocationResponse)
	if err == nil {
		return geolocationResponse, nil
	} else if err != domain.ErrDataNotFound {
		s.lg.WithError(err).Info("failed to retrieve data from cache")
		return nil, &domain.HTTPError{
			Message:    "failed to retrieve data from cache",
			StatusCode: http.StatusInternalServerError,
		}
	}

	geolocation, err := s.opts.Repo.GetGeolocation(ip)
	if err != nil {
		s.lg.WithError(err).Info("failed to retrieve geolocation data")
		return nil, &domain.HTTPError{
			Message:    "failed to retrieve geolocation data",
			StatusCode: http.StatusInternalServerError,
		}
	}

	go func() {
		s.opts.RedisRepo.Set(ctx, ip, geolocation)
		if err := s.opts.PlanRepo.UpdateUserCredits(userID, userPlan.CreditsUsed+1); err != nil {
			s.lg.WithError(err).Info("failed to update user credits")
		}
	}()

	return geolocation, nil
}
