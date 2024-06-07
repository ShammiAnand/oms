package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Secret      string
	TokenExpiry int
}

var Envs = initConfig()

func initConfig() Config {

	godotenv.Load()

	tokenExpiry := getEnv("TOKEN_EXPIRY", 3600).(int)

	return Config{
		Secret:      getEnv("SECRET", "some_secret").(string),
		TokenExpiry: tokenExpiry,
	}
}

func getEnv(key string, fallback interface{}) interface{} {
	if value, ok := os.LookupEnv(key); ok {
		switch fallback.(type) {
		case string:
			return value
		case int:
			if intValue, err := strconv.Atoi(value); err == nil {
				return intValue
			}
		case float64:
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				return floatValue
			}
		default:
			return value
		}
	}
	return fallback
}
