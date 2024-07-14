package connection

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection() {

	dsn := "Sahil:password@tcp(127.0.0.1:3306)/Sahil?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Can-not Connected to Database")
	} else {
		fmt.Println("Connected to DB...")
	}
}
