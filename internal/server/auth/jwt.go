package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var jwtSecret []byte

func init() {
	_ = godotenv.Load() // Загружаем .env, если есть, но не ругаемся, если нет

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("⚠️ Внимание: JWT_SECRET не найден в .env, проверь переменные окружения!")
		secret = "default_secret" // Установи значение по умолчанию или выйди с ошибкой
	}

	jwtSecret = []byte(secret)
}

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("невалидный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {

		return 0, errors.New("неверный формат claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return 0, errors.New("срок действия токена истёк")
	}

	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return 0, errors.New("неверный формат userID")
	}
	userID := int(userIDFloat)

	return userID, nil
}
