package main

import (
	"cloud/internal/config"
	"cloud/internal/server"
	"cloud/internal/server/storage"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	err := storage.InitDB()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}

	router := server.SetupServer(cfg)
	server.StartServer(router)

}
