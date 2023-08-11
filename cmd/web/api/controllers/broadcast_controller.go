package controllers

import (
	vmodels "customer_engagement/cmd/web/api/view_models"
	service "customer_engagement/service"
	"fmt"

	"encoding/json"

	"net/http"
)

type BroadcastController struct {
	service *service.Service
}

func NewBroadCastController(service *service.Service) *BroadcastController {
	return &BroadcastController{
		service: service,
	}
}

func (c BroadcastController) BroadcastGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var broadCastRequest vmodels.BroadcastRequest
		json.NewDecoder(r.Body).Decode(&broadCastRequest)
		ok, errors := broadCastRequest.Validate()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}

		_, err := c.service.Broadcast.EnqueueBroadcastSimpleSmsToGroup(r.Context(), broadCastRequest.MessageBody, broadCastRequest.GroupId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			fmt.Println("SUCCEDDED")
		}
	}
}
