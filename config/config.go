package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Type     string
		Sqlite   SqliteConfig
		Postgres PostgresConfig
	}
}

type SqliteConfig struct {
	Path string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
