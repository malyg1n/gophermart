package handler

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gophermart/pkg/config"
	"gophermart/pkg/token"
	v1 "gophermart/service/v1"
	"gophermart/storage"
	"gophermart/storage/pgsql"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Suite struct {
	suite.Suite
	handler *Handler
	us      storage.IUserStorage
	os      storage.IOrderStorage
	ts      storage.ITransactionStorage
}

func (s *Suite) SetupTest() {
	cfg, _ := config.GetConfig()
	cfg.DatabaseURI = "postgres://forge:secret@localhost:54321/gophermart?sslmode=disable"
	st, _ := pgsql.NewStorage(cfg)
	st.Truncate()
	us := v1.NewUserService(st, st)
	os := v1.NewOrderService(st, st)
	s.handler = NewHandler(WithUserService(us), WithOrderService(os))
	s.us = st
	s.os = st
	s.ts = st
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) testRequest(method, path string, payload io.Reader, headers map[string]string) *http.Response {
	ts := httptest.NewServer(s.handler.GetRouter())
	req, err := http.NewRequest(method, ts.URL+path, payload)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	require.NoError(s.T(), err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(s.T(), err)

	return resp
}

func (s *Suite) getAuthHeader(userID int) map[string]string {
	tkn, err := token.CreateTokenByUserID(userID)
	require.NoError(s.T(), err)

	return map[string]string{
		"Authorization": "Bearer " + tkn,
	}
}
