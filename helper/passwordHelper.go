package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedpassword), nil
}

func ComparePassword(hashedpassword, plainpassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(plainpassword))
}
