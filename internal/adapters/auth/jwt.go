package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/pkg/logger"

	"time"

	"github.com/dgrijalva/jwt-go"
)

// Auth is an authentication service that handles token creation and verification.
type Auth struct {
	secret string
	lg     *logger.Entry
}

// NewAuth creates a new instance of the Auth service.
func NewAuth(secret string) *Auth {
	return &Auth{
		secret: secret,
		lg:     logger.General.Component("Auth repository"),
	}
}

// Create generates a new JWT token with the provided user information.
func (a *Auth) Create(data domain.UserClaims) (*domain.JWTToken, error) {
	expirationTime := time.Now().Add(3 * time.Hour)
	claims := &JWTClaim{
		UserID: data.UserID,
		Email:  data.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.secret))
	if err != nil {
		a.lg.WithError(err).Error("failed to sign JWT token")
		return nil, fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return &domain.JWTToken{Token: signedToken, ExpirationTime: expirationTime}, nil
}

// Verify validates the provided JWT token and returns the user information.
func (a *Auth) Verify(signedToken string) (*domain.UserClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			a.lg.Error("invalid signing method used in JWT token")
			return nil, errors.New("invalid signing method used in JWT token")
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		a.lg.WithError(err).Error("failed to parse JWT token")
		return nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		a.lg.Error("invalid JWT token")
		return nil, errors.New("invalid JWT token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		a.lg.Error("JWT token has expired")
		return nil, errors.New("JWT token has expired")
	}

	return &domain.UserClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}

func (a *Auth) Authorize(ctx context.Context, token string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Auth) GetUserClaims(ctx context.Context, token string) (*domain.UserClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Auth) GetUserClaimsByEmail(ctx context.Context, email string) (*domain.UserClaims, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Auth) GetUser(ctx context.Context, token string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

// JWTClaim represents the claims in a JWT token.
type JWTClaim struct {
	Email  string `json:"email"`
	UserID string `json:"userID"`
	jwt.StandardClaims
}
