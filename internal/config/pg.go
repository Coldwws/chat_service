package config

import "fmt"

type PGConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

func loadPG() PGConfig {
	return PGConfig{
		Host:     mustEnv("PG_HOST"),
		Port:     getEnv("PG_PORT", "5432"),
		DBName:   mustEnv("PG_DATABASE_NAME"),
		User:     mustEnv("PG_USER"),
		Password: mustEnv("PG_PASSWORD"),
		SSLMode:  getEnv("PG_SSL_MODE", "disable"),
	}
}

func (c PGConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host,c.Port, c.DBName, c.User, c.Password, c.SSLMode,
	)
}