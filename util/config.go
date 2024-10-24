package util

import "github.com/spf13/viper"

type Config struct {
	DBURL      string `mapstructure:"DB_URL"`
	ServerAddr string `mapstructure:"SERVER_ADDR"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
