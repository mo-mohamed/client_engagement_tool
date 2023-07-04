package dbRepository

import (
	"testing"
	"time"

	testH "customer_engagement/test_helper"

	"gopkg.in/go-playground/assert.v1"
)

var (
	profileRepo IProfileRepo = NewProfileRepo(testH.DB)
)

func TestProfileDBRepository(t *testing.T) {
	t.Run("successfully create a profile", func(t *testing.T) {
		testH.TruncateTables([]string{"profile"})
		profile := newProfile("first", "last", "123")
		dbProfile, err := profileRepo.CreateProfile(testH.Ctx, profile)
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
		testH.TruncateTables([]string{"profile"})
		profile := newProfile("first", "last", "123")
		dbProfile, err := profileRepo.CreateProfile(testH.Ctx, profile)
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
		dbProfile, err = profileRepo.UpdateProfile(testH.Ctx, dbProfile)
		assert.Equal(t, err, nil)
		assert.Equal(t, *dbProfile.FirstName, firstNameChange)
		assert.Equal(t, *dbProfile.LastName, lastNameChange)
		assert.Equal(t, dbProfile.MDN, mdnChange)

	})

	t.Run("successfully finds a profile by id", func(t *testing.T) {
		testH.TruncateTables([]string{"profile"})
		profile := newProfile("first", "last", "123")
		profileRepo.CreateProfile(testH.Ctx, profile)
		dbProfile, err := profileRepo.GetProfile(testH.Ctx, profile.ID)
		assert.Equal(t, err, nil)
		assert.NotEqual(t, dbProfile.ID, nil)

	})

	t.Run("returns an error if a profile not found", func(t *testing.T) {
		testH.TruncateTables([]string{"profile"})
		_, err := profileRepo.GetProfile(testH.Ctx, 2)
		assert.Equal(t, err.Error(), "record not found")
	})

	t.Run("get profiles paginated and ordered by created at field DESC", func(t *testing.T) {
		testH.TruncateTables([]string{"profile"})
		profileRepo.CreateProfile(testH.Ctx, newProfile("first1", "last1", "123542"))
		profileRepo.CreateProfile(testH.Ctx, newProfile("first2", "last2", "122343"))
		now := time.Now()
		time.Sleep(time.Second)
		profileRepo.CreateProfile(testH.Ctx, newProfile("first3", "last3", "1543523"))
		profileRepo.CreateProfile(testH.Ctx, newProfile("first4", "last4", "1333e23"))
		profilesList, err := profileRepo.GetProfilesPaginated(testH.Ctx, 2, 2)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(profilesList), 2)
		for _, p := range profilesList {
			if p.CreatedAt.Before(now) {
				t.Error("Created at should be larger than ", now)
			}
		}
	})
}

func newProfile(firstName, lastName, mdn string) *ProfileStore {
	return &ProfileStore{
		FirstName: &firstName,
		LastName:  &lastName,
		MDN:       mdn,
	}
}
