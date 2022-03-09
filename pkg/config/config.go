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

func init() {
	viper.AutomaticEnv()
	pflag.StringP(dbURIName, "d", "", "database connection string")
	pflag.StringP(RunAddrName, "a", "localhost:8080", "run address")
	pflag.StringP(AccrualAddrName, "r", "http://localhost:8081", "accrual system address")
	pflag.StringP(AppSecretName, "k", "app-secret-key", "secret key for app")
}

// NewDefaultConfig return an instance of AppConfig.
func NewDefaultConfig() (config AppConfig, err error) {
	pflag.Parse()

	config = AppConfig{}
	err = viper.BindPFlags(pflag.CommandLine)

	if err != nil {
		return
	}

	config.RunAddress = viper.GetString(RunAddrName)
	config.DatabaseURI = viper.GetString(dbURIName)
	config.AccrualAddress = viper.GetString(AccrualAddrName)
	config.AppSecret = viper.GetString(AppSecretName)

	if config.RunAddress == "" {
		config.RunAddress = "localhost:8080"
	}
	if config.DatabaseURI == "" {
		return config, fmt.Errorf("database dsn was not setted")
	}
	if config.AccrualAddress == "" {
		return config, fmt.Errorf("accrual system adress was not setted")
	}
	if config.AppSecret == "" {
		config.AppSecret = "very-secret-key"
	}
	return
}
