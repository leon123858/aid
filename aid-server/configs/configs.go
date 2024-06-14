package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"
)

var Configs = struct {
	Host struct {
		Host string
		Port string
	}
	Jwt struct {
		Secret string
		time.Duration
	}
}{}

func init() {
	Configs.Host.Host = "localhost"
	Configs.Host.Port = "8080"
	Configs.Jwt.Secret = "test"
	Configs.Jwt.Duration = time.Hour
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	Configs.Host.Host = os.Getenv("HOST")
	Configs.Host.Port = os.Getenv("PORT")
	Configs.Jwt.Secret = os.Getenv("JWT_SECRET")
	Configs.Jwt.Duration, err = time.ParseDuration(os.Getenv("JWT_DURATION"))
}
