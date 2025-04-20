package models

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectedDB() {
	connectionString := "host=localhost user=postgres password=Sempakbau12 dbname=DbUser port=5432"

	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		panic("error connect database")
	}

	var migrationerr []error
	migrationerr = append(migrationerr, db.AutoMigrate(&UserAddressDetails{}))

	migrationerr = append(migrationerr, db.AutoMigrate(&DataNewUser{}))

	for _, migerr := range migrationerr {
		if migerr != nil {
			panic(fmt.Sprintf("failed connescted database %s ", migerr))
		}
	}
	DB = db
}
