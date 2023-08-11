package controllers

import (
	vmodels "customer_engagement/cmd/web/api/view_models"
	service "customer_engagement/service"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ProfileController struct {
	service *service.Service
}

func NewProfileController(ps *service.Service) *ProfileController {
	return &ProfileController{
		service: ps,
	}
}

func (c ProfileController) InitializeRoutes(r *mux.Router) {
	r.HandleFunc("/profile/create", c.Create()).Methods("POST")
	r.HandleFunc("/group/profile/add", c.AddToGroup()).Methods("POST")
}

func (c ProfileController) AddToGroup() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var groupProfileVm vmodels.GroupProfile
		err := json.NewDecoder(r.Body).Decode(&groupProfileVm)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.service.Profile.AttachToGroup(r.Context(), groupProfileVm.ProfileId, groupProfileVm.GroupId)
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
		json.NewDecoder(r.Body).Decode(&profileVM)

		ok, errors := profileVM.Validate()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}

		domianProfile := profileVM.ToDomain()
		profile, err := pc.service.Profile.Create(r.Context(), &domianProfile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		profileVM = profileVM.FromService(*profile)

		jsonResponse, err := json.Marshal(profileVM)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}

}
