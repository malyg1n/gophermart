package handler

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"gophermart/pkg/token"
	"gophermart/provider/accrual"
	orderService "gophermart/service/order/v1"
	userService "gophermart/service/user/v1"
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
	us      storage.UserStorer
	os      storage.OrderStorer
	ts      storage.TransactionStorer
}

func (s *Suite) SetupTest() {
	cfg, err := config.NewDefaultConfig()
	require.NoError(s.T(), err)
	cfg.DatabaseURI = "postgres://forge:secret@localhost:54321/gophermart?sslmode=disable"
	st, _ := pgsql.NewStorage(cfg.DatabaseURI)
	st.Truncate()

	lgr := logger.NewDefaultLogger()
	accrualProvider := accrual.NewFakeHTTProvider()

	us := userService.NewUserService(
		userService.WithUserStorageUserOption(st),
		userService.WithTransactionStorageUserOption(st),
		userService.WithLoggerUserOption(lgr),
	)

	ose := orderService.NewOrderService(
		orderService.WithOrderStorageOrderOption(st),
		orderService.WithTransactionStorageOrderOption(st),
		orderService.WithLoggerOrderOption(lgr),
		orderService.WithProviderOrderOption(accrualProvider),
	)

	s.handler = NewHandler(WithUserService(us), WithOrderService(ose), WithLogger(lgr))
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

func (s *Suite) getAuthHeader(userID uint64) map[string]string {
	tkn, err := token.CreateTokenByUserID(userID)
	require.NoError(s.T(), err)

	return map[string]string{
		"Authorization": "Bearer " + tkn,
	}
}
