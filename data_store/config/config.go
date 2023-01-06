package dbconfig

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	dbHost := os.Getenv("DBHOST")
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASSWORD")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")
	dsn := dbUser + ":" + dbPassword + "@tcp" + "(" + dbHost + ":" + dbPort + ")/" + dbName + "?" + "parseTime=true&loc=Local"
	fmt.Println(dsn)
	db_conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database : error=%v", err)
		os.Exit(1)
	}
	fmt.Println("Initializing database completed")
	DB = db_conn
}
