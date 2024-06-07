package utils

import (
	"errors"
	"github.com/behrouz-rfa/gateway-service/internal/core/common"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetAuthToken(gCtx *gin.Context) (string, error) {
	authHeader := gCtx.Request.Header.Get(string(common.AuthorizationContextKey))

	if !strings.Contains(authHeader, "Bearer") {
		return "", errors.New("invalid auth header")
	}

	const bearerHeader = "Bearer "

	if len(authHeader) < len(bearerHeader) {
		return "", errors.New("invalid auth header")
	}

	token := authHeader[len(bearerHeader):]

	if token == "" {
		return "", errors.New("empty bearer header")
	}

	return token, nil
}
