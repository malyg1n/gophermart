package handler

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"gophermart/pkg/token"
	"gophermart/provider/accrual"
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

	lgr := logger.GetLogger()
	accrualProvider := accrual.NewFakeHTTProvider()

	us := v1.NewUserService(
		v1.WithUserStorageUserOption(st),
		v1.WithTransactionStorageUserOption(st),
		v1.WithLoggerUserOption(lgr),
	)

	os := v1.NewOrderService(
		v1.WithOrderStorageOrderOption(st),
		v1.WithTransactionStorageOrderOption(st),
		v1.WithLoggerOrderOption(lgr),
		v1.WithProviderOrderOption(accrualProvider),
	)

	s.handler = NewHandler(WithUserService(us), WithOrderService(os), WithLogger(lgr))
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
