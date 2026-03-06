package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// TODO: think of a alternative solution
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local"
	}
	envPath := fmt.Sprintf(".env.%s", env)
	err := godotenv.Load(envPath)
	if err != nil {
		err = godotenv.Load(fmt.Sprintf("./user-service/%s", envPath))
		if err != nil {
			return err
		}
	}
	return nil
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
