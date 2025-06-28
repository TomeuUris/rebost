package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the database and returns a connection.
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("rebost.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&InventoryItem{}, &Product{}, &Nutriments{}, &Ration{})

	return db, nil
}
