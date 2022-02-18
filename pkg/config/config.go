package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AppConfig struct.
type AppConfig struct {
	RunAddress     string
	DatabaseURI    string
	AccrualAddress string
}

const (
	dbURIName       = "database_uri"
	RunAddrName     = "run_address"
	AccrualAddrName = "accrual_system_address"
)

var instance *AppConfig

// GetConfig return an instance of AppConfig.
func GetConfig() (*AppConfig, error) {
	if instance != nil {
		return instance, nil
	}

	viper.AutomaticEnv()
	pflag.StringP(dbURIName, "d", "", "database connection string")
	pflag.StringP(RunAddrName, "a", "http://0.0.0.0:8080", "run address")
	pflag.StringP(AccrualAddrName, "r", "http://127.0.0.1:8080", "accrual system address")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, fmt.Errorf("parsing flag error %w", err)
	}

	dbURI := viper.GetString(dbURIName)
	runAddr := viper.GetString("run_address")
	accrualAddr := viper.GetString("accrual_system_address")

	instance = &AppConfig{
		RunAddress:     runAddr,
		DatabaseURI:    dbURI,
		AccrualAddress: accrualAddr,
	}

	return instance, nil
}
