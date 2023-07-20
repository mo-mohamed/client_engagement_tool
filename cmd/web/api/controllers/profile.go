package controllers

import (
	vmodels "customer_engagement/cmd/web/api/view_models"
	service "customer_engagement/service"

	"encoding/json"
	"net/http"
)

type ProfileController struct {
	profileService service.IProfileService
}

func NewProfileController(ps service.IProfileService) *ProfileController {
	return &ProfileController{
		profileService: ps,
	}
}

func (c ProfileController) AddToGroup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var groupProfileVm vmodels.GroupProfile
		err := json.NewDecoder(r.Body).Decode(&groupProfileVm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.profileService.AttachToGroup(r.Context(), groupProfileVm.ProfileId, groupProfileVm.GroupId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

func (pc ProfileController) Create() func(w http.ResponseWriter, r *http.Request) {

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

		databaseProfile := profileVM.ToDatabaseEntity()

		db_profile, err := pc.profileService.Create(r.Context(), &databaseProfile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		profileVM = profileVM.FromDatabaseEntity(*db_profile)
		jsonResponse, err := json.Marshal(profileVM)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}

}
