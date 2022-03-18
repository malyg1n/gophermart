package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophermart/pkg/token"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
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

func (s *Suite) TestHandlerWithdrawals() {
	err := s.handler.userService.Create(context.Background(), "test", "test")
	require.NoError(s.T(), err)
	tkn, err := s.handler.userService.Auth(context.Background(), "test", "test")
	require.NoError(s.T(), err)
	userClaims, err := token.GetClaimsByToken(tkn)
	require.NoError(s.T(), err)

	err = s.handler.orderService.CreateOrder(context.Background(), "12345678903001", userClaims.UserID)
	require.NoError(s.T(), err)
	time.Sleep(time.Millisecond * 500)

	err = s.handler.userService.Withdraw(context.Background(), userClaims.UserID, "123456", 200)
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

func (s *Suite) TestHandlerWithdraw() {
	err := s.handler.userService.Create(context.Background(), "test", "test")
	require.NoError(s.T(), err)
	tkn, err := s.handler.userService.Auth(context.Background(), "test", "test")
	require.NoError(s.T(), err)
	userClaims, err := token.GetClaimsByToken(tkn)
	require.NoError(s.T(), err)

	err = s.handler.orderService.CreateOrder(context.Background(), "12345678903001", userClaims.UserID)
	require.NoError(s.T(), err)
	time.Sleep(time.Millisecond * 500)

	h1 := map[string]string{
		"Authorization": "Bearer " + tkn,
	}

	tests := []struct {
		name           string
		exceptedStatus int
		headers        map[string]string
		payload        io.Reader
	}{
		{
			name:           "200",
			exceptedStatus: 200,
			headers:        h1,
			payload:        strings.NewReader(`{"order": "12345678903", "sum": 300}`),
		},
		{
			name:           "401",
			exceptedStatus: 401,
			headers:        map[string]string{},
			payload:        strings.NewReader(`{"order": "12345678903", "sum": 300}`),
		},
		{
			name:           "402",
			exceptedStatus: 402,
			headers:        h1,
			payload:        strings.NewReader(`{"order": "12345678903", "sum": 300}`),
		},
		{
			name:           "422",
			exceptedStatus: 402,
			headers:        h1,
			payload:        strings.NewReader(`{"order": "", "sum": 300}`),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := s.testRequest(http.MethodPost, "/api/user/balance/withdraw", tt.payload, tt.headers)
			defer r.Body.Close()
			assert.Equal(s.T(), tt.exceptedStatus, r.StatusCode)
		})
	}
}
