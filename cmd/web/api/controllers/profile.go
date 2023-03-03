package controllers

import (
	viewModels "customer_engagement/cmd/web/api/view_models"
	dbconfig "customer_engagement/data_store/config"
	db_models "customer_engagement/data_store/models"
	repository "customer_engagement/data_store/repository"
	"strconv"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ProfileController struct{}

func (ProfileController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var profileVM viewModels.Profile
		err := json.NewDecoder(r.Body).Decode(&profileVM)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbProfile := profileVM.ToDTO()

		pRepo := repository.NewRepository[db_models.Profile](dbconfig.DB)
		err = pRepo.Add(&dbProfile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		profileVM = profileVM.FromDTO(dbProfile)
		jsonResponse, err := json.Marshal(profileVM)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}

}

func (ProfileController) AllByGroup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		group_id, err := strconv.Atoi(vars["group_id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var profileVM viewModels.Profile
		profiles := make([]viewModels.Profile, 0)
		pRepo := repository.NewRepository[db_models.Profile](dbconfig.DB)
		dbprofiles, err := pRepo.Where(&db_models.Profile{GroupID: &group_id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, v := range dbprofiles {
			profiles = append(profiles, profileVM.FromDTO(v))
		}

		jsonResponse, err := json.Marshal(profiles)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}

}
