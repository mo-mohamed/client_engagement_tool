package main

import (
	"customer_engagement/cmd/web/api/controllers"
	"customer_engagement/consumers"
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

var (
	dbConn *gorm.DB
)

func main() {
	setupDB()
	startConsumers()

	storeLayer := initStore()
	serviceLayer := initService(storeLayer)

	profileController := controllers.NewProfileController(serviceLayer)
	groupController := controllers.NewGroupController(serviceLayer)
	broadcastController := controllers.NewBroadCastController(serviceLayer)

	router := mux.NewRouter()
	profileController.InitializeRoutes(router)
	groupController.InitializeRoutes(router)
	broadcastController.InitializeRoutes(router)
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
	dbConn = db_conn
}

func startConsumers() {
	fmt.Println("Starting consumers")
	go consumers.NewGroupQueueConsumer().Run()
}

func initStore() *storeLayer.Store {
	return &storeLayer.Store{
		Profile: storeRepository.NewProfileRepo(dbConn),
		Group:   storeRepository.NewGroupRepo(dbConn),
	}
}

func initService(store *storeLayer.Store) *service.Service {
	return &service.Service{
		Profile: service.NewProfileService(store),
		Group:   service.NewGroupService(store),
	}
}
