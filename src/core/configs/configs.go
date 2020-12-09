package configs

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func Env(key string) string {
	return os.Getenv(key)
}

func IsDebugMode() bool {
	return Env("APP_DEBUG") == "true"
}
