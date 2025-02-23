package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
