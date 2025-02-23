package storage

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB инициализирует базу данных и создает таблицу пользователей
func InitDB(dataSourceName string) error {
	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite3", dataSourceName)
		if err != nil {
			return
		}

		// Проверяем соединение
		if err = db.Ping(); err != nil {
			return
		}

		// Создаём таблицу пользователей, если её нет
		query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL
		);`
		_, err = db.Exec(query)
	})

	return err
}

// GetDB возвращает подключение к базе данных
func GetDB() *sql.DB {
	return db
}

// GetUserByUsername ищет пользователя в БД по username
func GetUserByUsername(db *sql.DB, username string) (int, string, error) {
	var id int
	var passwordHash string
	query := `SELECT id, password_hash FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&id, &passwordHash)
	if err != nil {
		return 0, "", err
	}
	return id, passwordHash, nil
}

func CreateUser(db *sql.DB, username string, passwordHash string) error {
	_, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	return err
}
