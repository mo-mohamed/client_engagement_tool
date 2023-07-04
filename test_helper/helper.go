package testHelper

import (
	"context"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB = OpenDatabase()
var Ctx context.Context = context.TODO()

func OpenDatabase() *gorm.DB {
	dsn := os.Getenv("DB_TEST_USER") + ":" + os.Getenv("DB_TEST_PASSWORD") + "@tcp" +
		"(" + os.Getenv("DB_TEST_HOST") + ":" + os.Getenv("DB_TEST_PORT") + ")/" +
		os.Getenv("DB_TEST_NAME") + "?" + "parseTime=true&loc=Local"

	db_conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		os.Exit(1)
	}
	return db_conn
}

func TruncateTables(tables []string) {
	DB.Exec("SET FOREIGN_KEY_CHECKS=0;")
	for _, v := range tables {
		DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", v))
	}
	DB.Exec("SET FOREIGN_KEY_CHECKS=1;")

}
