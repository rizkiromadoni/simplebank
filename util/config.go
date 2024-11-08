package util

import (
	"flag"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env                  string        `mapstructure:"ENV"`
	AllowedOrigins       []string      `mapstructure:"ALLOWED_ORIGINS"`
	DBURL                string        `mapstructure:"DB_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddr       string        `mapstructure:"HTTP_SERVER_ADDR"`
	GRPCServerAddr       string        `mapstructure:"GRPC_SERVER_ADDR"`
	TokenSymetricKey     string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)

	if flag.Lookup("test.v") == nil {
		viper.SetConfigName("app")
		viper.SetConfigType("env")
	} else {
		viper.SetConfigName("test")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
