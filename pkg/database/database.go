package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
	//burda bazı kısımlar yapılacak olabilir
	fmt.Println("Database'e bağlantı gerçekleşti.")
	return db
}
