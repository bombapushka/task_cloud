package handlers

import (
	"cloud/internal/server/auth"
	"cloud/internal/server/storage"
	"cloud/internal/server/utils"
	"html/template"
	"log"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/auth.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Println("Ошибка загрузки шаблона:", err)
		return
	}

	if r.Method == http.MethodGet {
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	//db := storage.GetDB()

	userID, passwordHash, err := storage.GetUserByUsername(username)
	if err != nil || !utils.CheckPasswordHash(password, passwordHash) {
		data := map[string]interface{}{
			"Error": "Неверное имя пользователя или пароль",
		}
		tmpl.Execute(w, data)
		return
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		data := map[string]interface{}{
			"Error": "Ошибка генерации токена",
		}
		tmpl.Execute(w, data)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true, // Защита от XSS
		Path:     "/",
	})

	log.Println("Пользователь удачно вошел")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method, r.URL.Path)

	tmpl, err := template.ParseFiles("templates/auth.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Println("Ошибка загрузки шаблона:", err)
		return
	}

	if r.Method == http.MethodGet {
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" && password == "" {
		log.Println("Пользователь ввел пустые данные при регистрации")
		data := map[string]interface{}{
			"Error": "Данные должны быть заполнены",
		}
		tmpl.Execute(w, data)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Ошибка хэширования пароля:", err)
		data := map[string]interface{}{
			"Error": "Ошибка пароля",
		}
		tmpl.Execute(w, data)
		return
	}

	//db := storage.GetDB()

	if err := storage.CreateUser(username, hashedPassword); err != nil {
		log.Println("Пользователь уже существует:", err)
		data := map[string]interface{}{
			"Error": "Пользователь уже существует",
		}
		tmpl.Execute(w, data)
		return
	}

	userID, _, err := storage.GetUserByUsername(username)
	if err != nil {
		data := map[string]interface{}{
			"Error": "Ошибка в имени пользователя или пароли",
		}
		tmpl.Execute(w, data)
		return

	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		log.Println("Ошибка генерации токена:", err)
		data := map[string]interface{}{
			"Error": "Ошибка генерации токена",
		}
		tmpl.Execute(w, data)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})

	log.Println("Пользователь удачно зарегистрировался")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",             // Важно, чтобы совпадал с тем, что было при установке
		Expires:  time.Unix(0, 0), // Дата в прошлом
		MaxAge:   -1,              // Немедленное удаление
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
