package bcrypt

import (
	"errors"
	b "golang.org/x/crypto/bcrypt"
)

var ErrorPasswordTooLong = errors.New("Password too long")

func Encrypt(password string, cost int) (string, error) {
	bytePass := []byte(password)
	if len(bytePass) > 72 {
		return "", ErrorPasswordTooLong
	}
	hashed, err := b.GenerateFromPassword(bytePass, cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func Compare(original, hashed string) (bool, error) {
	err := b.CompareHashAndPassword([]byte(hashed), []byte(original))
	if err != nil {
		return false, err
	}
	return true, nil
}
