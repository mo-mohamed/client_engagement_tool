package controllers

import (
	dbconfig "customer_engagement/data_store/config"
	db_models "customer_engagement/data_store/models"
	repository "customer_engagement/data_store/repository"
	"encoding/json"
	"strconv"

	"customer_engagement/service/producers"
	"net/http"
)

type BroadcastController struct{}

func (BroadcastController) BroadcastGroup() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var bodyParams map[string]string
		json.NewDecoder(r.Body).Decode(&bodyParams)

		groupId, _ := strconv.Atoi(bodyParams["group_id"])
		gRepo := repository.NewRepository[db_models.Group](dbconfig.DB)
		dbGroup, err := gRepo.GetById(groupId)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
			return
		}

		producer := producers.GroupMessageProducer{}
		producer.EnqueueGroupBroadcast(dbGroup.ID, bodyParams["message_body"])

	}
}
