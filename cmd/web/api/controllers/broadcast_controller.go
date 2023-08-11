package controllers

import (
	vmodels "customer_engagement/cmd/web/api/view_models"
	service "customer_engagement/service"
	"fmt"

	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
)

type BroadcastController struct {
	service *service.Service
}

func NewBroadCastController(service *service.Service) *BroadcastController {
	return &BroadcastController{
		service: service,
	}
}

func (c BroadcastController) InitializeRoutes(r *mux.Router) {
	r.HandleFunc("/broadcast/sms", c.BroadcastGroup()).Methods("POST")
}

func (c BroadcastController) BroadcastGroup() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var bcr vmodels.BroadcastRequest
		json.NewDecoder(r.Body).Decode(&bcr)
		ok, errors := bcr.Validate()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}

		_, err := c.service.Broadcast.EnqueueBroadcastSimpleSmsToGroup(r.Context(), bcr.MessageBody, bcr.GroupId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			fmt.Println("SUCCEDDED")
		}
	}
}
