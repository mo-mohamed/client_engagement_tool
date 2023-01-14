package main

import (
	"customer_engagement/cmd/web/api/controllers"
	dbconfig "customer_engagement/data_store/config"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	setupDB()
	var groupController controllers.GroupController
	var profileController controllers.ProfileController

	router := mux.NewRouter()
	router.HandleFunc("/groups", groupController.All()).Methods("GET")
	router.HandleFunc("/group/create", groupController.Create()).Methods("POST")
	router.HandleFunc("/group/deactivate/{id}", groupController.Deactivate()).Methods("POST")
	router.HandleFunc("/profile/create", profileController.Create()).Methods("POST")
	router.HandleFunc("/group/{group_id}/profiles", profileController.AllByGroup()).Methods("GET")
	http.Handle("/", router)

	http.ListenAndServe(":8080", router)

}

func setupDB() {
	dbHost := os.Getenv("DBHOST")
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASSWORD")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")
	dsn := dbUser + ":" + dbPassword + "@tcp" + "(" + dbHost + ":" + dbPort + ")/" + dbName + "?" + "parseTime=true&loc=Local"
	db_conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Initializing database completed")
	dbconfig.DB = db_conn
}
