package middleware

import (
	"github.com/behrouz-rfa/gateway-service/internal/adapters/http/handlers"
	"github.com/behrouz-rfa/gateway-service/internal/core/common"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/behrouz-rfa/gateway-service/pkg/utils"
	"github.com/gin-gonic/gin"
)

// authMiddleware is a middleware to check if the user is authenticated
func AuthMiddleware(token ports.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := utils.GetAuthToken(ctx)
		if err != nil {
			// Log error - will be reported to Sentry
			err := domain.ErrEmptyAuthorizationHeader
			handlers.HandleError(ctx, err)
			return
		}

		payload, err := token.Verify(accessToken)
		if err != nil {
			handlers.HandleError(ctx, err)
			return
		}

		ctx.Set(common.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
