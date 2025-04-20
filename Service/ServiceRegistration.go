package Service

import (
	"errors"
	"fmt"
	"test-case/helper"
	"test-case/models"

	"gorm.io/gorm"
)

// type AuthService struct {
// 	DB *gorm.DB
// }

// register create new user

func Login(email, password string) (*models.User, error) {
	var user models.User

	fmt.Println("email input :", email)
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println("error query database", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email or password invalid")
		}
		return nil, err
	}
	if err := helper.ComparePassword(user.Password, password); err != nil {
		return nil, errors.New("invalid email or pasword")
	}
	return &user, nil
}
