package controllers

import (
	"bytes"
	viewModels "customer_engagement/cmd/web/api/view_models"
	service "customer_engagement/service"
	storeLayer "customer_engagement/store"
	storeRepository "customer_engagement/store/repository"
	testH "customer_engagement/test_helper"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

type CreateGroupRequest struct {
	Name *string `json:"name"`
}

func TestGroupController(t *testing.T) {
	t.Run("HTTP - Create a group", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		store := &storeLayer.Store{
			Profile: storeRepository.NewProfileRepo(testH.DB),
			Group:   storeRepository.NewGroupRepo(testH.DB),
		}

		service := service.Service{
			Group:   service.NewGroupService(store),
			Profile: service.NewProfileService(store),
		}

		groupController := NewGroupController(&service)
		handler := http.HandlerFunc(groupController.Create())

		group_name := "group 1"
		group_payload, _ := json.Marshal(viewModels.Group{
			Name: &group_name,
		})

		req, err := http.NewRequest("POST", "/group/create", bytes.NewBuffer(group_payload))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		assert.Equal(t, recorder.Result().Status, "201 Created")

		var body_response viewModels.Group
		json.Unmarshal(recorder.Body.Bytes(), &body_response)

		assert.Equal(t, *body_response.Name, group_name)
		assert.Equal(t, body_response.DeletedAt, nil)

	})
}
