package server

import (
	"log"
	"net/http"
)

func StartServer(router http.Handler) {
	log.Println("Сервер запущен")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln("Ошибка запуска сервера")
	}
}
