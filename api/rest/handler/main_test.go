package handler

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gophermart/pkg/config"
	v1 "gophermart/service/v1"
	"gophermart/storage/pgsql"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Suite struct {
	suite.Suite
	handler *Handler
}

func (s *Suite) SetupTest() {
	cfg, _ := config.GetConfig()
	cfg.DatabaseURI = "postgres://forge:secret@localhost:54321/gophermart?sslmode=disable"
	st, _ := pgsql.NewStorage(cfg)
	st.Clear()
	us := v1.NewUserService(st)
	s.handler = NewHandler(WithUserService(us))
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