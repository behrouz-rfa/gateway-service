package geo

import (
	"encoding/json"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGeolocation(t *testing.T) {
	tests := []struct {
		name           string
		ip             string
		apiKey         string
		mockStatusCode int
		mockResponse   interface{}
		expectError    bool
		expected       *domain.GeolocationResponse
	}{
		{
			name:           "Successful geolocation retrieval",
			ip:             "8.8.8.8",
			apiKey:         "test-api-key",
			mockStatusCode: http.StatusOK,
			mockResponse: &domain.GeolocationResponse{
				CountryName: "United States",
				City:        domain.City{Name: "Mountain View"},
			},
			expectError: false,
			expected: &domain.GeolocationResponse{
				CountryName: "United States",
				City:        domain.City{Name: "Mountain View"},
			},
		},
		{
			name:           "HTTP request fails",
			ip:             "8.8.8.8",
			apiKey:         "test-api-key",
			mockStatusCode: http.StatusInternalServerError,
			mockResponse:   nil,
			expectError:    true,
			expected:       &domain.GeolocationResponse{},
		},
		{
			name:           "Non-OK response status",
			ip:             "8.8.8.8",
			apiKey:         "test-api-key",
			mockStatusCode: http.StatusBadRequest,
			mockResponse:   nil,
			expectError:    true,
			expected:       &domain.GeolocationResponse{},
		},
		{
			name:           "Invalid response body",
			ip:             "8.8.8.8",
			apiKey:         "test-api-key",
			mockStatusCode: http.StatusOK,
			mockResponse:   "invalid json",
			expectError:    true,
			expected:       &domain.GeolocationResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock HTTP server
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					json.NewEncoder(w).Encode(tt.mockResponse)
				} else {
					w.Write([]byte(tt.mockResponse.(string)))
				}
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			// Override the URL to point to the mock server
			s := &GeolocationServiceImpl{
				APIKey: tt.apiKey,
				Client: server.Client(),
			}

			// Call the method
			resp, err := s.GetGeolocation(tt.ip)

			// Assert results
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp)
			}
		})
	}
}
