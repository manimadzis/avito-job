package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerHost   string `yaml:"server_host" mapstructure:"server_host"`
	ServerPort   string `yaml:"server_port" mapstructure:"server_port"`
	DBHost       string `yaml:"db_host" mapstructure:"db_host"`
	DBPort       string `yaml:"db_port" mapstructure:"db_port"`
	DBUsername   string `yaml:"db_username" mapstructure:"db_username"`
	DBPassword   string `yaml:"db_password" mapstructure:"db_password"`
	DatabaseName string `yaml:"database_name" mapstructure:"database_name"`
	LogLevel     string `yaml:"log_level" mapstructure:"log_level"`
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
