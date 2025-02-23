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
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Ошибка загрузки .env файла:", err)
	}

	uploadsDir := os.Getenv("UPLOADS_DIR")

	return &Config{
		UploadsDir: uploadsDir,
	}
}
