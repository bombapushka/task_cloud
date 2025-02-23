package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	UploadsDir string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // Просто пробуем загрузить .env, но не паникуем

	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "uploads" // Значение по умолчанию
	}

	return &Config{
		UploadsDir: uploadsDir,
	}
}
