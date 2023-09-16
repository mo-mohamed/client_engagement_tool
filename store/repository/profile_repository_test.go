package repository

import (
	"testing"
	"time"

	"customer_engagement/store"
	testHelper "customer_engagement/test_helper"

	"gopkg.in/go-playground/assert.v1"
)

func TestProfileDBRepository(t *testing.T) {

	store := &store.Store{
		Profile: NewProfileRepo(testHelper.DB),
		Group:   NewGroupRepo(testHelper.DB),
	}

	t.Run("successfully create a profile", func(t *testing.T) {
		testHelper.TruncateTables([]string{"profile"})
		profile := newProfile("first", "last", "123")
		dbProfile, err := store.Profile.CreateProfile(testHelper.Ctx, profile)
		assert.Equal(t, err, nil)
		assert.Equal(t, *profile.FirstName, *dbProfile.FirstName)
		assert.Equal(t, *profile.LastName, *dbProfile.LastName)
		assert.Equal(t, profile.MDN, dbProfile.MDN)
		if dbProfile.ID <= 0 {
			t.Error("Id is not correctly assigned in the returned db object, got id: ", dbProfile.ID)
		}
		assert.NotEqual(t, profile.CreatedAt, nil)
		assert.NotEqual(t, profile.UpdatedAt, nil)

	})

	t.Run("successfully updates a profile", func(t *testing.T) {
		testHelper.TruncateTables([]string{"profile"})
		profile := newProfile("first", "last", "123")
		dbProfile, err := store.Profile.CreateProfile(testHelper.Ctx, profile)
		assert.Equal(t, err, nil)
		assert.Equal(t, *profile.FirstName, *dbProfile.FirstName)
		assert.Equal(t, *profile.LastName, *dbProfile.LastName)
		assert.Equal(t, profile.MDN, dbProfile.MDN)
		if dbProfile.ID <= 0 {
			t.Error("Id is not correctly assigned in the returned db object, got id: ", dbProfile.ID)
		}

		firstNameChange := "first1"
		lastNameChange := "last2"
		mdnChange := "456"
		dbProfile.FirstName = &firstNameChange
		dbProfile.LastName = &lastNameChange
		dbProfile.MDN = mdnChange
		dbProfile, err = store.Profile.UpdateProfile(testHelper.Ctx, dbProfile)
		assert.Equal(t, err, nil)
		assert.Equal(t, *dbProfile.FirstName, firstNameChange)
		assert.Equal(t, *dbProfile.LastName, lastNameChange)
		assert.Equal(t, dbProfile.MDN, mdnChange)

	})

	t.Run("successfully finds a profile by id", func(t *testing.T) {
		testHelper.TruncateTables([]string{"profile"})
		profile := newProfile("first", "last", "123")
		store.Profile.CreateProfile(testHelper.Ctx, profile)
		dbProfile, err := store.Profile.GetProfile(testHelper.Ctx, profile.ID)
		assert.Equal(t, err, nil)
		assert.NotEqual(t, dbProfile.ID, nil)

	})

	t.Run("returns an error if a profile not found", func(t *testing.T) {
		testHelper.TruncateTables([]string{"profile"})
		_, err := store.Profile.GetProfile(testHelper.Ctx, 2)
		assert.Equal(t, err.Error(), "record not found")
	})

	t.Run("get profiles paginated and ordered by created at field DESC", func(t *testing.T) {
		testHelper.TruncateTables([]string{"profile"})
		store.Profile.CreateProfile(testHelper.Ctx, newProfile("first1", "last1", "123542"))
		store.Profile.CreateProfile(testHelper.Ctx, newProfile("first2", "last2", "122343"))
		now := time.Now()
		time.Sleep(time.Second)
		store.Profile.CreateProfile(testHelper.Ctx, newProfile("first3", "last3", "1543523"))
		store.Profile.CreateProfile(testHelper.Ctx, newProfile("first4", "last4", "1333e23"))
		profilesList, err := store.Profile.GetProfilesPaginated(testHelper.Ctx, 2, 2)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(profilesList), 2)
		for _, p := range profilesList {
			if p.CreatedAt.Before(now) {
				t.Error("Created at should be larger than ", now)
			}
		}
	})

	t.Run("get profiles paginated per group and ordered by created at field DESC", func(t *testing.T) {
		testHelper.TruncateTables([]string{"profile", "group_profile", "group"})
		g, _ := store.Group.CreateGroup(testHelper.Ctx, newGroup("group"))

		p1, _ := store.Profile.CreateProfile(testHelper.Ctx, newProfile("first1", "last1", "123542"))
		store.Profile.AddProfileToGroup(testHelper.Ctx, p1.ID, g.ID)

		p2, _ := store.Profile.CreateProfile(testHelper.Ctx, newProfile("first2", "last2", "122343"))
		store.Profile.AddProfileToGroup(testHelper.Ctx, p2.ID, g.ID)

		g2, _ := store.Group.CreateGroup(testHelper.Ctx, newGroup("group 2"))

		p3, _ := store.Profile.CreateProfile(testHelper.Ctx, newProfile("first3", "last3", "1543523"))
		store.Profile.AddProfileToGroup(testHelper.Ctx, p3.ID, g2.ID)

		p4, _ := store.Profile.CreateProfile(testHelper.Ctx, newProfile("first4", "last4", "1333e23"))
		store.Profile.AddProfileToGroup(testHelper.Ctx, p4.ID, g2.ID)

		profilesList, err := store.Profile.GetGroupProfilesPaginated(testHelper.Ctx, g.ID, 0, 10, time.Now())
		assert.Equal(t, err, nil)
		assert.Equal(t, len(profilesList), 2)
	})
}
