package main

import (
	groupsHandler "customer_engagement/cmd/web/api/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/groups", groupsHandler.All).Methods("GET")
	router.HandleFunc("/group/create", groupsHandler.Create).Methods("POST")
	router.HandleFunc("/group/deactivate/{id}", groupsHandler.Deactivate).Methods("POST")
	http.Handle("/", router)

	http.ListenAndServe(":8080", router)

}
