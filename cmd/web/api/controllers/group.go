package controllers

import (
	vmodels "customer_engagement/cmd/web/api/view_models"

	service "customer_engagement/service"
	"strconv"
	"time"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type GroupController struct {
	service service.IGroupService
}

func NewGroupController(gs service.IGroupService) *GroupController {
	return &GroupController{
		service: gs,
	}
}

// func (c GroupController) All() func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		dbGroups, err := c.service.GetAll()
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		var vmGroup vmodels.Group
// 		vmGroups := make([]vmodels.Group, 0)
// 		for _, v := range *dbGroups {
// 			vmGroups = append(vmGroups, vmGroup.FromDatabaseEntity(v))
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		jsonResponse, err := json.Marshal(vmGroups)
// 		if err != nil {
// 			return
// 		}

// 		w.Write(jsonResponse)
// 	}
// }

func (c GroupController) Create() func(w http.ResponseWriter, r *http.Request) {
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
		dbEntity := group.ToDatabaseEntity()
		dbGroup, err := c.service.CreateGroup(r.Context(), &dbEntity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		group = group.FromDatabaseEntity(*dbGroup)
		jsonResponse, err := json.Marshal(group)
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}

func (c GroupController) Deactivate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Error occured", http.StatusBadRequest)
			return
		}

		dbGroup, err := c.service.GetGroup(r.Context(), groupID)
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
		_, err = c.service.UpdateGroup(r.Context(), dbGroup)
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
