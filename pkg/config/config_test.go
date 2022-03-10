package config

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetConfig(t *testing.T) {
	t.Setenv(strings.ToUpper(dbURIName), "fake-db-dsn")
	t.Setenv(strings.ToUpper(RunAddrName), "fake-addr")
	t.Setenv(strings.ToUpper(AccrualAddrName), "fake-accrual")
	got, err := NewDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, "fake-db-dsn", got.DatabaseURI)
	assert.Equal(t, "fake-addr", got.RunAddress)
	assert.Equal(t, "fake-accrual", got.AccrualAddress)
}
