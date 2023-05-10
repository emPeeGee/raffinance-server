package crypt

import "golang.org/x/crypto/bcrypt"

func HashPassphrase(passphrase string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passphrase), 14)
	return string(bytes), err
}

func CheckPasswordHashes(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
