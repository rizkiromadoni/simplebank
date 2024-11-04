package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBURL                string        `mapstructure:"DB_URL"`
	HTTPServerAddr       string        `mapstructure:"HTTP_SERVER_ADDR"`
	GRPCServerAddr       string        `mapstructure:"GRPC_SERVER_ADDR"`
	TokenSymetricKey     string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
