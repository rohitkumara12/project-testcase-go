package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primary key"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password" gorm:"unique"`
	Count     int       `gorm:"default:0"`
	LockUntil time.Time `gorm:"default:NULL"`
}
