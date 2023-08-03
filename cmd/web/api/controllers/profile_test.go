package controllers

import (
	service "customer_engagement/service"
	storeLayer "customer_engagement/store"
	storeRepository "customer_engagement/store/repository"
	testH "customer_engagement/test_helper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileService(t *testing.T) {
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

		req, err := http.NewRequest("POST", "/profile/create", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

	})
}
