package storage

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB инициализирует подключение к PostgreSQL
func InitDB() error {
	var err error
	once.Do(func() {
		dataSourceName := os.Getenv("DATABASE_URL") // Читаем из .env
		if dataSourceName == "" {
			log.Fatal("DATABASE_URL не установлен")
		}

		db, err = sql.Open("postgres", dataSourceName)
		if err != nil {
			log.Fatal("Ошибка при подключении к PostgreSQL:", err)
			return
		}

		// Проверяем соединение
		if err = db.Ping(); err != nil {
			log.Fatal("Не удалось подключиться к PostgreSQL:", err)
			return
		}

		log.Println("Успешное подключение к PostgreSQL")

		// Создаём таблицу пользователей, если её нет
		query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL
		);`
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal("Ошибка при создании таблицы:", err)
		}
	})

	return err
}

// GetDB возвращает подключение к базе данных
func GetDB() *sql.DB {
	return db
}

func GetUserByUsername(db *sql.DB, username string) (int, string, error) {
	var id int
	var passwordHash string
	query := `SELECT id, password_hash FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&id, &passwordHash)
	if err != nil {
		return 0, "", err
	}
	return id, passwordHash, nil
}

func CreateUser(db *sql.DB, username string, passwordHash string) error {
	_, err := db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, passwordHash)
	return err
}
