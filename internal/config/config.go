package config

import (
	"log"
	"os"
)

type Config struct {
	Env  string
	GRPC GRPCConfig
	PG   PGConfig
}

func LoadConfig() Config {
	return Config{
		Env:  getEnv("ENV","local"),
		GRPC: loadGRPC(),
		PG:   loadPG(),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == ""{
		log.Fatalf("missing env : %s",key)
	}
	return v
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == ""{
		return def
	}
	return v
}