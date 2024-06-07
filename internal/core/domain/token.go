package domain

import "time"

// UserClaims is an entity that represents the payload of the token
type UserClaims struct {
	UserID string
	Email  string
}

type JWTToken struct {
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expirationTime"`
}
