package ipservice

import (
	"errors"
	"net"
	"testing"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/stretchr/testify/assert"
)

// Mock client implementing the ipinfo.Client interface
type MockIPInfoClient struct {
	mockGetIPInfo func(ip net.IP) (*ipinfo.Core, error)
}

func (m *MockIPInfoClient) GetIPInfo(ip net.IP) (*ipinfo.Core, error) {
	return m.mockGetIPInfo(ip)
}

func TestGetIpInfo(t *testing.T) {
	tests := []struct {
		name        string
		ip          string
		mockResp    *ipinfo.Core
		mockErr     error
		expected    *ipinfo.Core
		expectError bool
	}{
		{
			name: "Successful IP info retrieval",
			ip:   "8.8.8.8",
			mockResp: &ipinfo.Core{
				IP:       net.ParseIP("8.8.8.8"),
				Hostname: "dns.google",
			},
			expected: &ipinfo.Core{
				IP:       net.ParseIP("8.8.8.8"),
				Hostname: "dns.google",
			},
			expectError: false,
		},
		{
			name:        "Invalid IP address",
			ip:          "invalid-ip",
			mockResp:    nil,
			mockErr:     errors.New("invalid IP address"),
			expected:    nil,
			expectError: true,
		},
		{
			name:        "API call fails",
			ip:          "8.8.8.8",
			mockResp:    nil,
			mockErr:     errors.New("API call failed"),
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockIPInfoClient{
				mockGetIPInfo: func(ip net.IP) (*ipinfo.Core, error) {
					return tt.mockResp, tt.mockErr
				},
			}

			repo := &IPServiceRepo{client: mockClient}

			resp, err := repo.GetIpInfo(tt.ip)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp)
			}
		})
	}
}
