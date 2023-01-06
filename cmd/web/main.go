package main

import (
	"customer_engagement/cmd/web/api/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
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
