package controllers

import (
	"bytes"
	viewModels "customer_engagement/cmd/web/api/view_models"
	service "customer_engagement/service"
	"customer_engagement/service/models"
	storeLayer "customer_engagement/store"
	storeRepository "customer_engagement/store/repository"
	testH "customer_engagement/test_helper"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

type AddProfileToGroup struct {
	ProfileId int `json:"profile_id"`
	GroupId   int `json:"group_id"`
}
type CreateProfile struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	MDN       *string `json:"mdn"`
}

func TestProfileController(t *testing.T) {
	t.Run("HTTP - Create a profile", func(t *testing.T) {
		testH.TruncateTables([]string{"`profile`"})
		store := &storeLayer.Store{
			Profile: storeRepository.NewProfileRepo(testH.DB),
			Group:   storeRepository.NewGroupRepo(testH.DB),
		}

		service := service.Service{
			Group:   service.NewGroupService(store),
			Profile: service.NewProfileService(store),
		}

		profileController := NewProfileController(&service)
		handler := http.HandlerFunc(profileController.Create())

		profile_first_name := "first"
		profile_last_name := "last"
		profile_mdn := "12345"
		profile_payload, _ := json.Marshal(CreateProfile{
			FirstName: &profile_first_name,
			LastName:  &profile_last_name,
			MDN:       &profile_mdn,
		})

		req, err := http.NewRequest("POST", "/api/v1/profile/create", bytes.NewBuffer(profile_payload))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		assert.Equal(t, recorder.Result().Status, "201 Created")

		var body_response viewModels.Profile
		json.Unmarshal(recorder.Body.Bytes(), &body_response)

		assert.Equal(t, *body_response.FirstName, "first")
		assert.Equal(t, *body_response.LastName, "last")
		assert.Equal(t, *body_response.MDN, "12345")
		assert.Equal(t, body_response.Active, true)
		assert.Equal(t, body_response.DeletedAt, nil)

	})

	t.Run("HTTP - Adds a profile to a group", func(t *testing.T) {
		testH.TruncateTables([]string{"profile", "group_profile", "`group`"})
		store := &storeLayer.Store{
			Profile: storeRepository.NewProfileRepo(testH.DB),
			Group:   storeRepository.NewGroupRepo(testH.DB),
		}

		service := service.Service{
			Group:   service.NewGroupService(store),
			Profile: service.NewProfileService(store),
		}

		group_name := "group name"
		group, err := service.Group.Create(testH.Ctx, &models.Group{Name: &group_name})
		assert.Equal(t, err, nil)

		profile_first_name := "first"
		profile_last_name := "last"
		profile_mdn := "12345"
		profile, err := service.Profile.Create(testH.Ctx, &models.Profile{
			FirstName: &profile_first_name,
			LastName:  &profile_last_name,
			MDN:       profile_mdn,
		})

		assert.Equal(t, err, nil)

		profileController := NewProfileController(&service)
		handler := http.HandlerFunc(profileController.AddToGroup())

		payload, _ := json.Marshal(AddProfileToGroup{
			ProfileId: profile.ID,
			GroupId:   group.ID,
		})

		req, err := http.NewRequest("POST", "/api/v1/group/profile/add", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		assert.Equal(t, recorder.Result().Status, "201 Created")
	})

	t.Run("HTTP - CREATE PROFILE - returns an error if MDN is missing", func(t *testing.T) {
		testH.TruncateTables([]string{"`profile`"})
		store := &storeLayer.Store{
			Profile: storeRepository.NewProfileRepo(testH.DB),
			Group:   storeRepository.NewGroupRepo(testH.DB),
		}

		service := service.Service{
			Group:   service.NewGroupService(store),
			Profile: service.NewProfileService(store),
		}

		profileController := NewProfileController(&service)
		handler := http.HandlerFunc(profileController.Create())

		profile_first_name := "first"
		profile_last_name := "last"
		profile_payload, _ := json.Marshal(CreateProfile{
			LastName:  &profile_last_name,
			FirstName: &profile_first_name,
		})

		req, err := http.NewRequest("POST", "/api/v1/profile/create", bytes.NewBuffer(profile_payload))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		assert.Equal(t, recorder.Result().StatusCode, http.StatusBadRequest)

		var body_response []viewModels.ValidationError
		json.Unmarshal(recorder.Body.Bytes(), &body_response)

		validation_error := body_response[0]

		assert.Equal(t, validation_error.Field, "MDN")
		assert.Equal(t, validation_error.Tag, "required")
	})
}
