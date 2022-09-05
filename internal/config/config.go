package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MongoConfig struct {
	URI      string
	Username string
	Password string
	Database string
}

type NetServerConfig struct {
	Host string
	Port int
}

type Config struct {
	Mongo  MongoConfig
	Server NetServerConfig
}

func parseConfigFile(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("serverListener.tcp", &cfg.Server); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *Config) error {
	if err := envconfig.Process("db", &cfg.Mongo); err != nil {
		return err
	}
	return nil
}

func Init(configDir string) (*Config, error) {
	viper.SetConfigName("config")
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := setFromEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
}
