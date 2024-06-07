package geo

import (
	"encoding/json"
	"fmt"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"

	"net/http"
)

type GeolocationServiceImpl struct {
	APIKey string
	Client *http.Client
}

func NewGeolocationService(apiKey string, client *http.Client) *GeolocationServiceImpl {
	if client == nil {
		client = http.DefaultClient
	}
	return &GeolocationServiceImpl{APIKey: apiKey, Client: client}
}

func (s *GeolocationServiceImpl) GetGeolocation(ip string) (*domain.GeolocationResponse, error) {
	response := new(domain.GeolocationResponse)

	url := fmt.Sprintf("https://api.ip2location.io/?key=%s&ip=%s", s.APIKey, ip)
	resp, err := s.Client.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("failed to get geolocation: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}

	return response, nil
}
