package handlers

import (
	"cloud/internal/config"
	"cloud/internal/server/auth"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Println("Ошибка загрузки шаблона:", err)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	// Получаем userID
	userID, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	userDir := filepath.Join(cfg.UploadsDir, fmt.Sprintf("%d", userID))

	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		if err := os.MkdirAll(userDir, 0755); err != nil {
			log.Println("Ошибка создания папки:", err)
			http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
			return
		}
	}

	// Получаем список файлов
	files, err := os.ReadDir(userDir)
	if err != nil {
		http.Error(w, "Ошибка чтения каталога", http.StatusInternalServerError)
		return
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	// Передаём данные в HTML-шаблон
	tmpl.Execute(w, struct {
		Files []string
	}{
		Files: fileList,
	})
}
