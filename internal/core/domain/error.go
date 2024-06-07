package domain

import (
	"net/http"
)

type HTTPError struct {
	Message    string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return e.Message
}

var (
	// ErrInternal is an error for when an internal service fails to process the request
	ErrInternal = &HTTPError{
		Message:    "internal error",
		StatusCode: http.StatusInternalServerError,
	}
	// ErrDataNotFound is an error for when requested data is not found
	ErrDataNotFound = &HTTPError{
		Message:    "data not found",
		StatusCode: http.StatusNotFound,
	}
	// ErrNoUpdatedData is an error for when no data is provided to update
	ErrNoUpdatedData = &HTTPError{
		Message:    "no data to update",
		StatusCode: http.StatusBadRequest,
	}
	// ErrConflictingData is an error for when data conflicts with existing data
	ErrConflictingData = &HTTPError{
		Message:    "data conflicts with existing data in unique column",
		StatusCode: http.StatusConflict,
	}
	// ErrInsufficientStock is an error for when product stock is not enough
	ErrInsufficientStock = &HTTPError{
		Message:    "product stock is not enough",
		StatusCode: http.StatusConflict,
	}
	// ErrInsufficientPayment is an error for when total paid is less than total price
	ErrInsufficientPayment = &HTTPError{
		Message:    "total paid is less than total credit",
		StatusCode: http.StatusBadRequest,
	}
	// ErrTokenDuration is an error for when the token duration format is invalid
	ErrTokenDuration = &HTTPError{
		Message:    "invalid token duration format",
		StatusCode: http.StatusBadRequest,
	}
	// ErrTokenCreation is an error for when the token creation fails
	ErrTokenCreation = &HTTPError{
		Message:    "error creating token",
		StatusCode: http.StatusInternalServerError,
	}
	// ErrExpiredToken is an error for when the access token is expired
	ErrExpiredToken = &HTTPError{
		Message:    "access token has expired",
		StatusCode: http.StatusUnauthorized,
	}
	// ErrInvalidToken is an error for when the access token is invalid
	ErrInvalidToken = &HTTPError{
		Message:    "access token is invalid",
		StatusCode: http.StatusUnauthorized,
	}
	// ErrInvalidCredentials is an error for when the credentials are invalid
	ErrInvalidCredentials = &HTTPError{
		Message:    "invalid email or password",
		StatusCode: http.StatusUnauthorized,
	}
	// ErrEmptyAuthorizationHeader is an error for when the authorization header is empty
	ErrEmptyAuthorizationHeader = &HTTPError{
		Message:    "authorization header is not provided",
		StatusCode: http.StatusUnauthorized,
	}
	// ErrInvalidAuthorizationHeader is an error for when the authorization header is invalid
	ErrInvalidAuthorizationHeader = &HTTPError{
		Message:    "authorization header format is invalid",
		StatusCode: http.StatusUnauthorized,
	}
	// ErrInvalidAuthorizationType is an error for when the authorization type is invalid
	ErrInvalidAuthorizationType = &HTTPError{
		Message:    "authorization type is not supported",
		StatusCode: http.StatusUnauthorized,
	}
	// ErrUnauthorized is an error for when the user is unauthorized
	ErrUnauthorized = &HTTPError{
		Message:    "user is unauthorized to access the resource",
		StatusCode: http.StatusUnauthorized,
	}
	// ErrForbidden is an error for when the user is forbidden to access the resource
	ErrForbidden = &HTTPError{
		Message:    "user is forbidden to access the resource",
		StatusCode: http.StatusForbidden,
	}
)
