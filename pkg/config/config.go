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
	AppSecret      string
}

const (
	dbURIName       = "database_uri"
	RunAddrName     = "run_address"
	AccrualAddrName = "accrual_system_address"
	AppSecretName   = "app_secret"
)

var instance *AppConfig

// GetConfig return an instance of AppConfig.
func GetConfig() (*AppConfig, error) {
	if instance != nil {
		return instance, nil
	}

	viper.AutomaticEnv()
	pflag.StringP(dbURIName, "d", "", "database connection string")
	pflag.StringP(RunAddrName, "a", "localhost:8080", "run address")
	pflag.StringP(AccrualAddrName, "r", "localhost:8081", "accrual system address")
	pflag.StringP(AppSecretName, "k", "app-secret-key", "secret key for app")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, fmt.Errorf("parsing flag error %w", err)
	}

	dbURI := viper.GetString(dbURIName)
	runAddr := viper.GetString(RunAddrName)
	accrualAddr := viper.GetString(AccrualAddrName)
	appSecret := viper.GetString(AppSecretName)

	if appSecret == "" {
		appSecret = "very-secret-key"
	}

	instance = &AppConfig{
		RunAddress:     runAddr,
		DatabaseURI:    dbURI,
		AccrualAddress: accrualAddr,
		AppSecret:      appSecret,
	}

	return instance, nil
}
