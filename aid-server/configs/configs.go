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
	Path struct {
		UserDB  string
		IDMap   string
		AliasDB string
	}
	Time struct {
		LoginCache time.Duration
	}
}{}

func init() {
	Configs.Host.Host = "0.0.0.0"
	Configs.Host.Port = "8080"
	Configs.Jwt.Secret = "test"
	Configs.Jwt.Duration = time.Minute * 1
	Configs.Path.UserDB = "data/user.db"
	Configs.Path.IDMap = "data/idmap.db"
	Configs.Path.AliasDB = "data/alias.db"
	Configs.Time.LoginCache = time.Minute * 1
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	Configs.Host.Host = os.Getenv("HOST")
	Configs.Host.Port = os.Getenv("PORT")
	Configs.Jwt.Secret = os.Getenv("JWT_SECRET")
	Configs.Jwt.Duration, _ = time.ParseDuration(os.Getenv("JWT_DURATION"))
	Configs.Path.UserDB = os.Getenv("USER_DB")
	Configs.Path.IDMap = os.Getenv("ID_MAP")
	Configs.Path.AliasDB = os.Getenv("ALIAS_DB")
	Configs.Time.LoginCache, _ = time.ParseDuration(os.Getenv("LOGIN_CACHE"))
}
