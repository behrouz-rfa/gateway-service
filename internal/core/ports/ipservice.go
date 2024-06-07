package ports

import (
	"context"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/ipinfo/go/v2/ipinfo"
)

type IpServiceRepository interface {
	GetIpInfo(ip string) (*ipinfo.Core, error)
}

type IpInfoService interface {
	GetIpInfo(ctx context.Context, ip, userID string) (*domain.IpInfo, error)
}
