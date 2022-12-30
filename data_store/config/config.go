package dbconfig

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	dsn := "mostafa" + ":" + "P@ssw0rd" + "@tcp" + "(" + "localhost" + ":" + "3306" + ")/" + "customer_eng" + "?" + "parseTime=true&loc=Local"
	db_conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database : error=%v", err)
		os.Exit(1)
	}
	fmt.Println("Initializing database completed")
	DB = db_conn
}
