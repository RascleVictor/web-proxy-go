package config

import (
	"log"

	"github.com/spf13/viper"
)

// ServerConfig contient la configuration du serveur et des backends
type ServerConfig struct {
	Port     int
	Backends []string
}

// LoggingConfig contient la configuration du logger
type LoggingConfig struct {
	Level string
}

// Config contient toute la configuration de l'application
type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
}

// LoadConfig charge le fichier YAML de config
func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return Config{
		Server: ServerConfig{
			Port:     viper.GetInt("server.port"),
			Backends: viper.GetStringSlice("server.backends"),
		},
		Logging: LoggingConfig{
			Level: viper.GetString("logging.level"),
		},
	}
}
