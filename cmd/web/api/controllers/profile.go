package controllers

import (
	vmodels "customer_engagement/cmd/web/api/view_models"
	dbc "customer_engagement/data_store/config"
	dbm "customer_engagement/data_store/models"
	repository "customer_engagement/data_store/repository"

	"encoding/json"
	"net/http"
)

type ProfileController struct{}

func (ProfileController) AddToGroup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var groupProfileVm vmodels.GroupProfile
		err := json.NewDecoder(r.Body).Decode(&groupProfileVm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbGroupProfile := groupProfileVm.ToDatabaseEntity()

		pRepo := repository.NewRepository[dbm.GroupProfile](dbc.DB)
		err = pRepo.Add(&dbGroupProfile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		groupProfileVm = groupProfileVm.FromDatabaseEntity(dbGroupProfile)
		jsonResponse, err := json.Marshal(groupProfileVm)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}

func (ProfileController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var profileVM vmodels.Profile
		err := json.NewDecoder(r.Body).Decode(&profileVM)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, errors := profileVM.Validate()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}

		dbProfile := profileVM.ToDatabaseEntity()

		pRepo := repository.NewRepository[dbm.Profile](dbc.DB)
		err = pRepo.Add(&dbProfile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		profileVM = profileVM.FromDatabaseEntity(dbProfile)
		jsonResponse, err := json.Marshal(profileVM)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}

}
