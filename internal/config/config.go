package config

import (
	"os"
)

func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}
