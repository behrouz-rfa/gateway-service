package handlers

import (
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IpInfoHandler struct {
	Service ports.IpInfoService
}

func NewIpInfoHandler(service ports.IpInfoService) *IpInfoHandler {
	return &IpInfoHandler{Service: service}
}

// GetIpInfo godoc
//
//	@Summary		Retrieve IP Information
//	@Description	Get detailed information for a specified IP address.
//	@Tags			IP
//	@Accept			json
//	@Produce		json
//	@Param			ip	path		string			true	"IP address"
//	@Success		200				{object}	domain.IpInfo	"IP information retrieved successfully"
//	@Failure		400				{object}	errorResponse	"Validation error: IP address is required"
//	@Failure		401				{object}	errorResponse	"Unauthorized error: Invalid or missing authentication token"
//	@Failure		404				{object}	errorResponse	"Data not found error: No information found for the given IP address"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error: Failed to retrieve IP information"
//	@Router			/ipinfo/{ip} [GET]
//	@Security		BearerAuth
func (h *IpInfoHandler) GetIpInfo(c *gin.Context) {
	payload := GetAuthPayload(c)

	ip := c.Param("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
		return
	}

	resp, err := h.Service.GetIpInfo(c, ip, payload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
