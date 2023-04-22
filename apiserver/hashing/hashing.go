package hashing

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// TODO: revise hashing cost
	// NOTE: GenerateFromPassword hashes and salts the password
	// think about using a pepper as well
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
