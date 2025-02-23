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

func UploadHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Ошибка разбора шаблона", http.StatusBadRequest)
		log.Println("Ошибка разбора шаблона:", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка загрузки файла", http.StatusBadRequest)
		log.Println("Ошибка загрузки файла", err)
		return
	}
	defer file.Close()

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	// Проверяем токен и получаем userID
	userID, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	// Создаём папку для пользователя
	uploadDir := filepath.Join(cfg.UploadsDir, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(uploadDir, fileHeader.Filename)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		http.Error(w, "Ошибка открытия файла", http.StatusBadRequest)
		log.Println("Ошибка открытия файла:", err)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		http.Error(w, "Ошибка копирования файла", http.StatusInternalServerError)
		log.Println("Ошибка копирования файла:", err)
		return
	}
	log.Println("Файл успешно загружен на сайта")

	fmt.Println("Загруженный файл:", fileHeader.Filename)
	fmt.Println("Размер файла:", fileHeader.Size)
	fmt.Println("MIME-тип:", fileHeader.Header.Get("Content-Type"))

	//w.Write([]byte("Файл успешно загружен"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
