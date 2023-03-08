package controllers

import (
	viewModels "customer_engagement/cmd/web/api/view_models"
	dbconfig "customer_engagement/data_store/config"
	db_models "customer_engagement/data_store/models"
	repository "customer_engagement/data_store/repository"

	"encoding/json"

	"customer_engagement/service/producers"
	"net/http"
)

type BroadcastController struct{}

func (BroadcastController) BroadcastGroup() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var bcr viewModels.BroadcastRequest
		json.NewDecoder(r.Body).Decode(&bcr)
		ok, errors := bcr.Validate()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}

		groupId := bcr.GroupId
		gRepo := repository.NewRepository[db_models.Group](dbconfig.DB)
		dbGroup, err := gRepo.GetById(groupId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		producer := producers.GroupMessageProducer{}
		err = producer.EnqueueGroupBroadcast(dbGroup.ID, bcr.MessageBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
}
