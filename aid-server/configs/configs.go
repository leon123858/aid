package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var Configs = struct {
	Host struct {
		Host string
		Port string
	}
}{}

func init() {
	Configs.Host.Host = "localhost"
	Configs.Host.Port = "8080"
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	Configs.Host.Host = os.Getenv("HOST")
	Configs.Host.Port = os.Getenv("PORT")
}
