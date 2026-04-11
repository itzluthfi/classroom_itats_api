package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=beruang3.itats.ac.id user=ak45 password=ationsio45 dbname=ak070203 port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var results []struct {
		UserName string `gorm:"column:name"`
		RoleName string `gorm:"column:role"`
	}

	err = db.Table("users").
		Select("users.name as name, role.name as role").
		Joins("JOIN users_roles ON users.uid = users_roles.uid").
		Joins("JOIN role ON role.rid = users_roles.rid").
		Where("users.name = ?", "412110170108").
		Scan(&results).Error

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Ditemukan %d role untuk user 412110170108:\n", len(results))
	for i, res := range results {
		fmt.Printf("%d. Role: %s\n", i+1, res.RoleName)
	}
}
