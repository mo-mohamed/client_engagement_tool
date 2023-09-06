package service

import (
	interfaces "customer_engagement/service/interfaces"
	domains "customer_engagement/service/models"
	storeLayer "customer_engagement/store"
	storeRepo "customer_engagement/store/repository"
	testH "customer_engagement/test_helper"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestProfileService(t *testing.T) {

	var storedb *storeLayer.Store = &storeLayer.Store{
		Profile: storeRepo.NewProfileRepo(testH.DB),
		Group:   storeRepo.NewGroupRepo(testH.DB),
	}
	var profileService interfaces.IProfileService = NewProfileService(storedb)

	t.Run("successfully create a profile", func(t *testing.T) {
		testH.TruncateTables([]string{"`profile`"})
		_, err := profileService.Create(testH.Ctx, newProfile("first", "last", "12343"))
		assert.Equal(t, err, nil)
	})

	t.Run("successfully gets a profile via id", func(t *testing.T) {
		testH.TruncateTables([]string{"`profile`"})
		p, _ := profileService.Create(testH.Ctx, newProfile("first", "last", "12343"))
		p, err := profileService.Get(testH.Ctx, p.ID)
		assert.Equal(t, err, nil)
		assert.Equal(t, *p.FirstName, "first")
	})

	t.Run("successfully updates a profile", func(t *testing.T) {
		testH.TruncateTables([]string{"`profile`"})
		p, _ := profileService.Create(testH.Ctx, newProfile("first", "last", "12343"))
		newFirstName := "firstname-updated"
		newLastName := "lastname-updated"
		p.FirstName = &newFirstName
		p.LastName = &newLastName
		g, err := profileService.Update(testH.Ctx, p)
		assert.Equal(t, err, nil)
		assert.Equal(t, *g.FirstName, newFirstName)
		assert.Equal(t, *g.LastName, newLastName)
	})
}

func newProfile(firstName, lastName, mdn string) *domains.Profile {
	return &domains.Profile{
		FirstName: &firstName,
		LastName:  &lastName,
		MDN:       mdn,
	}
}
