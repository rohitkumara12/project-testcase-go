package models

import (
	"gorm.io/gorm"
)

type UserAddressDetails struct {
	gorm.Model

	UserID  int    `json:"userid"`
	Street  string `json:"street"`
	City    string `gorm:"size:255" json:"city"`
	Country string `gorm:"size:255" json:"country"`
	Kodepos string `gorm:"size:255" json:"kodepos"`
}
