package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func (s *Suite) TestHandlerRegister() {
	tests := []struct {
		name           string
		payload        io.Reader
		exceptedStatus int
		exceptedHeader string
	}{
		{
			name:           "valid",
			payload:        strings.NewReader(`{"login": "test", "password": "secret"}`),
			exceptedStatus: 200,
			exceptedHeader: "Authorization",
		},
		{
			name:           "exists login",
			payload:        strings.NewReader(`{"login": "test", "password": "secret"}`),
			exceptedStatus: 409,
			exceptedHeader: "",
		},
		{
			name:           "empty login",
			payload:        strings.NewReader(`{"login": "", "password": "secret"}`),
			exceptedStatus: 400,
			exceptedHeader: "",
		},
	}
	headers := map[string]string{
		"Content-type": "application/json",
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := s.testRequest(http.MethodPost, "/api/user/register", tt.payload, headers)
			defer r.Body.Close()
			assert.Equal(s.T(), tt.exceptedStatus, r.StatusCode)
			if tt.exceptedHeader != "" {
				assert.NotEmpty(t, r.Header.Get(tt.exceptedHeader))
			}
		})
	}
}

func (s *Suite) TestHandlerLogin() {
	err := s.handler.userService.Create(context.Background(), "test", "secret")
	assert.NoError(s.T(), err)
	tests := []struct {
		name           string
		payload        io.Reader
		exceptedStatus int
		exceptedHeader string
	}{
		{
			name:           "valid",
			payload:        strings.NewReader(`{"login": "test", "password": "secret"}`),
			exceptedStatus: 200,
			exceptedHeader: "Authorization",
		},
		{
			name:           "invalid login",
			payload:        strings.NewReader(`{"login": "test1", "password": "secret"}`),
			exceptedStatus: 401,
			exceptedHeader: "",
		},
		{
			name:           "empty login",
			payload:        strings.NewReader(`{"login": "", "password": "secret"}`),
			exceptedStatus: 400,
			exceptedHeader: "",
		},
	}
	headers := map[string]string{
		"Content-type": "application/json",
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := s.testRequest(http.MethodPost, "/api/user/login", tt.payload, headers)
			defer r.Body.Close()
			assert.Equal(s.T(), tt.exceptedStatus, r.StatusCode)
			if tt.exceptedHeader != "" {
				assert.NotEmpty(t, r.Header.Get(tt.exceptedHeader))
			}
		})
	}
}
