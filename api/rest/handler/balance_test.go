package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophermart/pkg/token"
	"net/http"
	"testing"
)

func (s *Suite) TestHandlerBalance() {
	err := s.handler.userService.Create(context.Background(), "test", "test")
	assert.NoError(s.T(), err)
	tkn, err := s.handler.userService.Auth(context.Background(), "test", "test")
	assert.NoError(s.T(), err)
	h1 := map[string]string{
		"Authorization": "Bearer " + tkn,
	}

	tests := []struct {
		name           string
		exceptedStatus int
		headers        map[string]string
	}{
		{
			name:           "200",
			exceptedStatus: 200,
			headers:        h1,
		},
		{
			name:           "401",
			exceptedStatus: 401,
			headers:        map[string]string{},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := s.testRequest(http.MethodGet, "/api/user/balance", nil, tt.headers)
			defer r.Body.Close()
			assert.Equal(t, tt.exceptedStatus, r.StatusCode)
		})
	}
}

func (s *Suite) TestHandlerWithdraws() {
	err := s.handler.userService.Create(context.Background(), "test", "test")
	require.NoError(s.T(), err)
	tkn, err := s.handler.userService.Auth(context.Background(), "test", "test")
	require.NoError(s.T(), err)
	userID, err := token.GetUserIDByToken(tkn)
	require.NoError(s.T(), err)

	err = s.handler.userService.TopUp(context.Background(), userID, "123456", 500)
	require.NoError(s.T(), err)
	err = s.handler.userService.Withdraw(context.Background(), userID, "123456", 200)
	require.NoError(s.T(), err)

	h1 := map[string]string{
		"Authorization": "Bearer " + tkn,
	}

	err = s.handler.userService.Create(context.Background(), "test1", "test1")
	require.NoError(s.T(), err)
	tkn, err = s.handler.userService.Auth(context.Background(), "test1", "test1")
	require.NoError(s.T(), err)

	h2 := map[string]string{
		"Authorization": "Bearer " + tkn,
	}

	tests := []struct {
		name           string
		exceptedStatus int
		headers        map[string]string
	}{
		{
			name:           "200",
			exceptedStatus: 200,
			headers:        h1,
		},
		{
			name:           "204",
			exceptedStatus: 204,
			headers:        h2,
		},
		{
			name:           "401",
			exceptedStatus: 401,
			headers:        map[string]string{},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := s.testRequest(http.MethodGet, "/api/user/balance/withdrawals", nil, tt.headers)
			defer r.Body.Close()
			assert.Equal(s.T(), tt.exceptedStatus, r.StatusCode)
		})
	}
}
