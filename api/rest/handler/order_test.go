package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func (s *Suite) TestHandlerCreateOrder() {
	h1 := s.getAuthHeader(100)
	h2 := s.getAuthHeader(101)

	tests := []struct {
		name           string
		exceptedStatus int
		headers        map[string]string
		payload        io.Reader
	}{
		{
			name:           "401",
			exceptedStatus: 401,
			headers:        map[string]string{},
			payload:        strings.NewReader(""),
		},
		{
			name:           "202",
			exceptedStatus: 202,
			headers:        h1,
			payload:        strings.NewReader("12345678903"),
		},
		{
			name:           "200",
			exceptedStatus: 200,
			headers:        h1,
			payload:        strings.NewReader("12345678903"),
		},
		{
			name:           "422",
			exceptedStatus: 422,
			headers:        h1,
			payload:        strings.NewReader(""),
		},
		{
			name:           "409",
			exceptedStatus: 409,
			headers:        h2,
			payload:        strings.NewReader("12345678903"),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := s.testRequest(http.MethodPost, "/api/user/orders", tt.payload, tt.headers)
			defer r.Body.Close()
			assert.Equal(t, tt.exceptedStatus, r.StatusCode)
		})
	}
}

func (s *Suite) TestHandlerGetOrdersByUser() {
	h1 := s.getAuthHeader(100)
	h2 := s.getAuthHeader(101)
	s.handler.orderService.CreateOrder(context.Background(), "12345678903", 100)

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
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			r := s.testRequest(http.MethodGet, "/api/user/orders", nil, tt.headers)
			defer r.Body.Close()
			assert.Equal(t, tt.exceptedStatus, r.StatusCode)
		})
	}
}
