package services

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"
	"net/http"
)

// IpInfoServiceOpts holds the configuration for IpInfoService.
type IpInfoServiceOpts struct {
	IpRepo    ports.IpServiceRepository
	PlanRepo  ports.PlanRepository
	RedisRepo ports.RedisRepository
}

// IpInfoService provides services related to IP information.
type IpInfoService struct {
	opts IpInfoServiceOpts
	lg   *logger.Entry
}

// NewIpInfoService creates a new instance of IpInfoService.
func NewIpInfoService(opts IpInfoServiceOpts) *IpInfoService {
	return &IpInfoService{
		opts: opts,
		lg:   logger.General.Component("IpInfoService"),
	}
}

// GetIpInfo retrieves IP information for the given IP address and user ID.
func (s *IpInfoService) GetIpInfo(ctx context.Context, ip string, userID string) (*domain.IpInfo, error) {
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

	geolocationResponse := new(domain.IpInfo)
	err = s.opts.RedisRepo.Get(ctx, ip, geolocationResponse)
	if err == nil {
		return geolocationResponse, nil
	}

	coreInfo, err := s.opts.IpRepo.GetIpInfo(ip)
	if err != nil {
		s.lg.WithError(err).Info("failed to retrieve IP information")
		return nil, &domain.HTTPError{
			Message:    "failed to retrieve IP information",
			StatusCode: http.StatusInternalServerError,
		}
	}

	info := &domain.IpInfo{
		IP:             coreInfo.IP.String(),
		Hostname:       coreInfo.Hostname,
		Bogon:          coreInfo.Bogon,
		Anycast:        coreInfo.Anycast,
		City:           coreInfo.City,
		Region:         coreInfo.Region,
		Country:        coreInfo.Country,
		CountryFlagURL: coreInfo.CountryFlagURL,
		IsEU:           coreInfo.IsEU,
		Location:       coreInfo.Location,
		Org:            coreInfo.Org,
		Postal:         coreInfo.Postal,
		Timezone:       coreInfo.Timezone,
	}

	go func() {
		s.opts.RedisRepo.Set(ctx, ip, info)

		if err := s.opts.PlanRepo.UpdateUserCredits(userID, userPlan.CreditsUsed+1); err != nil {
			s.lg.WithError(err).Info("failed to update user credits")
		}
	}()

	return info, nil
}
