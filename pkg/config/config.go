package config

import (
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

func init() {
	viper.AutomaticEnv()
	pflag.StringP(dbURIName, "d", "", "database connection string")
	pflag.StringP(RunAddrName, "a", "localhost:8080", "run address")
	pflag.StringP(AccrualAddrName, "r", "http://localhost:8081", "accrual system address")
	pflag.StringP(AppSecretName, "k", "app-secret-key", "secret key for app")
}

// NewDefaultConfig return an instance of AppConfig.
func NewDefaultConfig() (*AppConfig, error) {
	if instance != nil {
		return instance, nil
	}

	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)

	if err != nil {
		return nil, err
	}

	instance := &AppConfig{}
	instance.RunAddress = viper.GetString(RunAddrName)
	instance.DatabaseURI = viper.GetString(dbURIName)
	instance.AccrualAddress = viper.GetString(AccrualAddrName)
	instance.AppSecret = viper.GetString(AppSecretName)

	return instance, nil
}
