// Package config gerencia as configurações da aplicação.
// Carrega variáveis de ambiente do arquivo .env e fornece
// uma struct tipada para acesso seguro às configurações.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	MongoUri          string
	MongoRootUsername string
	MongoRootPassword string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		MongoUri:          getEnv("MONGO_URI", "mongodb://root:example@localhost:27017/?authSource=admin"),
		MongoRootUsername: getEnv("MONGO_ROOT_USERNAME", "root"),
		MongoRootPassword: getEnv("MONGO_ROOT_PASSWORD", "example"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
