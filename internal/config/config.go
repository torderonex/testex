package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string           `yaml:"env"`
	Os         string           `yaml:"os"`
	Postgres   PostgresDatabase `yaml:"postgres"`
	HTTPServer HTTPServer       `mapstructure:"http_server"`
}

type HTTPServer struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresDatabase struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
		log.Fatalf("config file doesn't exists: %s", configPath)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to decode config into struct: %v", err)
	}

	return config
}
