package handlers

import (
	"cloud/internal/config"
	"cloud/internal/server/auth"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		http.Error(w, "Ошибка валидации токена", http.StatusUnauthorized)
		return
	}

	fileName := r.URL.Query().Get("filename")
	userDir := filepath.Join(cfg.UploadsDir, fmt.Sprintf("%d", userID))
	filePath := filepath.Join(userDir, fileName)

	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Файл не найден", http.StatusNotFound)
			log.Println("Файл не найден", err)
			return
		} else {
			http.Error(w, "Ошибка доступа", http.StatusNotFound)
			log.Println("Ошибка доступа:", err)
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Ошибка открытия файла", http.StatusNotFound)
		log.Println("Ошибка открытия файла")
		return
	}
	defer file.Close()

	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	log.Println("Началось скачивание файла:", fileName)
	io.Copy(w, file)
}
