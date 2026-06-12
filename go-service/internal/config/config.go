package config

import (
	"os"
)

type Config struct {
	Port     string
	Root     string
	Password string
	Dbname   string
	Host     string
	Token    string
}

func NewConfig() *Config {

	return &Config{
		Port:     os.Getenv("DB_PORT"),
		Root:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Dbname:   os.Getenv("POSTGRES_DB"),
		Host:     os.Getenv("HOST"),
		Token:    os.Getenv("BOT_TOKEN"),
	}
}
