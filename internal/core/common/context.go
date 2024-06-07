package common

type ContextKey string

const (
	AuthorizationContextKey ContextKey = "Authorization"
	GinContextKey           ContextKey = "gin"
	UserClaimContextKey     ContextKey = "user_claim"
	UserContextKey          ContextKey = "user"
	AuthorizationHeaderKey             = "authorization"
	AuthorizationType                  = "bearer"
	AuthorizationPayloadKey            = "authorization_payload"
)
