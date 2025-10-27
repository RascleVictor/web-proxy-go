package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port   int
	Target string
}

type LoggingConfig struct {
	Level string
}

type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return Config{
		Server: ServerConfig{
			Port:   viper.GetInt("server.port"),
			Target: viper.GetString("server.target"),
		},
		Logging: LoggingConfig{
			Level: viper.GetString("logging.level"),
		},
	}
}
