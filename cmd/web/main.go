package main

import (
	"customer_engagement/cmd/web/api/controllers"
	"customer_engagement/consumers"
	queueService "customer_engagement/queue"
	awsSqsClient "customer_engagement/queue/awssqs"
	service "customer_engagement/service"
	storeLayer "customer_engagement/store"
	storeRepository "customer_engagement/store/repository"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	queueClient := initQueueService()
	serviceLayer := initService(storeLayer, queueClient)
	// serviceLayer.Queue = queueClient

	profileController := controllers.NewProfileController(serviceLayer)
	groupController := controllers.NewGroupController(serviceLayer)
	broadcastController := controllers.NewBroadCastController(serviceLayer)

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	api1 := api.PathPrefix("/v1").Subrouter()
	api1.HandleFunc("/group/create", groupController.Create()).Methods("POST")
	api1.HandleFunc("/group/deactivate/{id}", groupController.Deactivate()).Methods("POST")

	api1.HandleFunc("/profile/create", profileController.Create()).Methods("POST")
	api1.HandleFunc("/group/profile/add", profileController.AddToGroup()).Methods("POST")

	api1.HandleFunc("/broadcast/sms", broadcastController.BroadcastGroup()).Methods("POST")

	http.Handle("/", api)

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

func initService(store *storeLayer.Store, queueClient queueService.IQueueClient) *service.Service {
	return &service.Service{
		Profile:   service.NewProfileService(store),
		Group:     service.NewGroupService(store),
		Broadcast: service.NewBroadcastService(store, queueClient),
	}
}

func initQueueService() queueService.IQueueClient {
	sqsSession := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_SQS_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_SQS_ENDPOINT")),
	}))

	return awsSqsClient.NewSqs(sqs.New(sqsSession))
}
