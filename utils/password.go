package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing input passwd")
	}

	return string(hashedPasswd), nil
}

func ComparePassword(password string, hashedPasswd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(password))
}
