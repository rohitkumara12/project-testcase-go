package models

import (
	"time"

	"gorm.io/gorm"
)

type DataNewUser struct {
	gorm.Model
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password" gorm:"unique"`
	Count     int       `gorm:"default:0"`
	LockUntil time.Time `gorm:"default:NULL"`
}
