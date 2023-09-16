package repository

import (
	"testing"

	"customer_engagement/store"
	"customer_engagement/store/models"
	testHelper "customer_engagement/test_helper"

	"gopkg.in/go-playground/assert.v1"
)

func TestMTRepository(t *testing.T) {
	store := &store.Store{
		Profile: NewProfileRepo(testHelper.DB),
		Group:   NewGroupRepo(testHelper.DB),
		MT:      NewMTRepository(testHelper.DB),
	}

	t.Run("successfully creates a mobile terminated record", func(t *testing.T) {
		testHelper.TruncateTables([]string{"mobile_terminated", "profile", "`group`"})
		profile := newProfile("first name", "last_name", "12345")
		_, err := store.Profile.CreateProfile(testHelper.Ctx, profile)
		assert.Equal(t, err, nil)

		group := newGroup("gorup name")
		_, err = store.Group.CreateGroup(testHelper.Ctx, group)
		assert.Equal(t, err, nil)

		mt := models.MT{
			ProfileId:   profile.ID,
			GroupId:     group.ID,
			BroadcastId: "broadcast id",
		}

		dbMt, err := store.MT.CreateMT(testHelper.Ctx, &mt)
		assert.Equal(t, err, nil)
		assert.Equal(t, dbMt.Processed, false)
		assert.Equal(t, dbMt.BroadcastId, mt.BroadcastId)
		assert.Equal(t, dbMt.ProfileId, profile.ID)
		assert.Equal(t, dbMt.GroupId, group.ID)
	})

	t.Run("Checks if MT is processed", func(t *testing.T) {
		testHelper.TruncateTables([]string{"mobile_terminated", "profile", "`group`"})
		profile := newProfile("first name", "last_name", "12345")
		_, err := store.Profile.CreateProfile(testHelper.Ctx, profile)
		assert.Equal(t, err, nil)

		group := newGroup("gorup name")
		_, err = store.Group.CreateGroup(testHelper.Ctx, group)
		assert.Equal(t, err, nil)

		mt := &models.MT{
			ProfileId:   profile.ID,
			GroupId:     group.ID,
			BroadcastId: "broadcast id",
		}

		processed := store.MT.IsProcessed(testHelper.Ctx, mt)
		assert.Equal(t, processed, false)

	})
}
