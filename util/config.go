package util

import "github.com/spf13/viper"

type Config struct {
	DatabaseDriver string `mapstructure:"DB_DRIVER"`
	DatabaseSource string `mapstructure:"DB_SOURCE"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
}

func LoadEnv(path string) (config Config, err error) {
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
