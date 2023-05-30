package controllers

import (
	vmodels "customer_engagement/cmd/web/api/view_models"
	dbc "customer_engagement/data_store/config"
	dbm "customer_engagement/data_store/models"
	repository "customer_engagement/data_store/repository"
	"strconv"
	"time"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type GroupController struct{}

func (GroupController) All() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gRepo := repository.NewRepository[dbm.Group](dbc.DB)
		dbGroups, err := gRepo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var vmGroup vmodels.Group
		vmGroups := make([]vmodels.Group, 0)
		for _, v := range *dbGroups {
			vmGroups = append(vmGroups, vmGroup.FromDatabaseEntity(v))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse, err := json.Marshal(vmGroups)
		if err != nil {
			return
		}

		w.Write(jsonResponse)
	}
}

func (GroupController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var group vmodels.Group
		err := json.NewDecoder(r.Body).Decode(&group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, errors := group.Validate()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}

		dbGroup := group.ToDatabaseEntity()
		gRepo := repository.NewRepository[dbm.Group](dbc.DB)
		err = gRepo.Add(&dbGroup)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		group = group.FromDatabaseEntity(dbGroup)
		jsonResponse, err := json.Marshal(group)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}

func (GroupController) Deactivate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
			return
		}

		gRepo := repository.NewRepository[dbm.Group](dbc.DB)
		dbGroup, err := gRepo.GetById(groupID)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
			return
		}

		if dbGroup.DeletedAt != nil {
			http.Error(w, "Group alreade deactivated", http.StatusBadRequest)
			return
		}
		deleted_at := time.Now()
		dbGroup.DeletedAt = &deleted_at
		err = gRepo.Update(dbGroup)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
			return
		}

		var vmGroup vmodels.Group
		g := vmGroup.FromDatabaseEntity(*dbGroup)
		jsonResponse, err := json.Marshal(g)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}

}
