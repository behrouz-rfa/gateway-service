package ports

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
)

type GeolocationRepository interface {
	GetGeolocation(ip string) (*domain.GeolocationResponse, error)
}

type GeolocationService interface {
	GetGeolocation(ctx context.Context, ip string, userID string) (*domain.GeolocationResponse, error)
}
