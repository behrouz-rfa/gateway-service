package ipservice

import (
	"github.com/ipinfo/go/v2/ipinfo"
	"net"
)

// Define an interface for the IPInfo client
type IPInfoClient interface {
	GetIPInfo(ip net.IP) (*ipinfo.Core, error)
}

type IPServiceRepo struct {
	client IPInfoClient
}

func NewIPServiceRepo(apiKey string) *IPServiceRepo {
	client := ipinfo.NewClient(nil, nil, apiKey)
	return &IPServiceRepo{client: client}
}

func (i *IPServiceRepo) GetIpInfo(ip string) (*ipinfo.Core, error) {

	info, err := i.client.GetIPInfo(net.ParseIP(ip))

	if err != nil {
		return nil, err
	}

	return info, nil
}
