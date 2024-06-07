package handlers

import (
	"github.com/behrouz-rfa/gateway-service/internal/core/common"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GeolocationHandler struct {
	Service ports.GeolocationService
}

func NewGeolocationHandler(service ports.GeolocationService) *GeolocationHandler {
	return &GeolocationHandler{Service: service}
}

// GetGeolocation godoc
//
//	@Summary		Get Geolocation Information
//	@Description	Retrieve geolocation information for a given IP address.
//	@Tags			GEO
//	@Accept			json
//	@Produce		json
//	@Param			ip	path		string			true	"IP address"
//	@Success		200				{object}	userResponse	"Geolocation information retrieved successfully"
//	@Failure		400				{object}	errorResponse	"Validation error: IP address is required"
//	@Failure		401				{object}	errorResponse	"Unauthorized error: Invalid or missing authentication token"
//	@Failure		404				{object}	errorResponse	"Data not found error: No geolocation information found for the given IP address"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error: Failed to retrieve geolocation information"
//	@Router			/geo/{ip} [GET]
//	@Security		BearerAuth
func (h *GeolocationHandler) GetGeolocation(c *gin.Context) {
	payload := GetAuthPayload(c)

	ip := c.Param("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
		return
	}

	response, err := h.Service.GetGeolocation(c, ip, payload.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAuthPayload is a helper function to get the auth payload from the context
func GetAuthPayload(ctx *gin.Context) *domain.UserClaims {
	return ctx.MustGet(common.AuthorizationPayloadKey).(*domain.UserClaims)
}
