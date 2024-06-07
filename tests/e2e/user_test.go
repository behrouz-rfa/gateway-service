//go:build e2e

package gql

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/behrouz-rfa/gateway-service/internal/core/common"
	"github.com/goccy/go-json"
)

func (s *GqlTestSuite) TestCreateAndGetUser() {
	// Register User
	inputReq := `{
		"name": "Jon",
		"password": "12345678",
		"email":"joa2@gmail.com"
	}`
	req, err := http.NewRequest("POST", "/api/v1/users/register", strings.NewReader(inputReq))
	if err != nil {
		s.T().Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.api.ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var userData userResponseData
	err = json.Unmarshal(w.Body.Bytes(), &userData)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(userData.Data.Email, "joa2@gmail.com")

	// Get User
	req2, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", userData.Data.ID), nil)
	if err != nil {
		s.T().Fatal(err)
	}
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set(string(common.AuthorizationContextKey), "Bearer "+userData.Data.JWTToken.Token)

	w = httptest.NewRecorder()
	s.api.ServeHTTP(w, req2)
	s.Equal(http.StatusOK, w.Code)

	var fetchedUserData userResponseData
	err = json.Unmarshal(w.Body.Bytes(), &fetchedUserData)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(fetchedUserData.Data.Name, "Jon")
}
