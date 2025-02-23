package config

import (
	"github.com/joho/godotenv"
	"log"
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

	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadsDir, 0755); err != nil {
			log.Fatal("Ошибка создания директории uploads: " + err.Error())
		}
	}

	return &Config{
		UploadsDir: uploadsDir,
	}
}
