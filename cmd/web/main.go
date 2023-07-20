package main

import (
	"customer_engagement/cmd/web/api/controllers"
	"customer_engagement/consumers"
	dbconfig "customer_engagement/data_store/config"
	service "customer_engagement/service"
	storeLayer "customer_engagement/store"
	storeRepository "customer_engagement/store/repository"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	setupDB()
	startConsumers()

	storeLayer := &storeLayer.Store{
		Profile: storeRepository.NewProfileRepo(dbconfig.DB),
		Group:   storeRepository.NewGroupRepo(dbconfig.DB),
	}
	profileService := service.NewProfileService(storeLayer)
	groupService := service.NewGroupService(storeLayer)

	profileController := controllers.NewProfileController(profileService)
	groupController := controllers.NewGroupController(groupService)
	broadcastController := controllers.NewBroadCastController(groupService)
	// var groupController controllers.GroupController

	router := mux.NewRouter()
	// router.HandleFunc("/groups", groupController.All()).Methods("GET")
	router.HandleFunc("/group/create", groupController.Create()).Methods("POST")
	router.HandleFunc("/group/deactivate/{id}", groupController.Deactivate()).Methods("POST")
	router.HandleFunc("/profile/create", profileController.Create()).Methods("POST")
	router.HandleFunc("/broadcast/sms", broadcastController.BroadcastGroup()).Methods("POST")
	router.HandleFunc("/group/profile/add", profileController.AddToGroup()).Methods("POST")
	http.Handle("/", router)

	http.ListenAndServe(":8080", router)

}

func setupDB() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := dbUser + ":" + dbPassword + "@tcp" + "(" + dbHost + ":" + dbPort + ")/" + dbName + "?" + "parseTime=true&loc=Local"
	fmt.Println("Initializing database connection")
	db_conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Successful connection to database")
	dbconfig.DB = db_conn
}

func startConsumers() {
	fmt.Println("Starting consumers")
	go consumers.NewGroupQueueConsumer().Run()
}
