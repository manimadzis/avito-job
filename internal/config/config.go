package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerHost          string `mapstructure:"server_host"`
	ServerPort          string `mapstructure:"server_port"`
	DBHost              string `mapstructure:"db_host"`
	DBPort              string `mapstructure:"db_port"`
	DBUsername          string `mapstructure:"db_username"`
	DBPassword          string `mapstructure:"db_password"`
	DatabaseName        string `mapstructure:"database_name"`
	LogLevel            string `mapstructure:"log_level"`
	FileServerDirectory string `mapstructure:"file_server_directory"`
}

func Load(src string) (*Config, error) {
	viper.SetConfigFile(src)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, newErrCantLoadConfig(err)
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, newErrCantParseConfig(err)
	}

	return &config, nil
}
