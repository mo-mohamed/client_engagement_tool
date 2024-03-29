package repository

import (
	"customer_engagement/store"
	testHelper "customer_engagement/test_helper"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestGroupRepository(t *testing.T) {
	store := &store.Store{
		Profile: NewProfileRepo(testHelper.DB),
		Group:   NewGroupRepo(testHelper.DB),
	}

	t.Run("successfully inserts a group", func(t *testing.T) {
		testHelper.TruncateTables([]string{"`group`"})
		groupName := "Group 1"
		dbGroup, err := store.Group.CreateGroup(testHelper.Ctx, newGroup(groupName))
		assert.Equal(t, err, nil)
		assert.Equal(t, dbGroup.Name, groupName)
		assert.NotEqual(t, dbGroup.CreatedAt, nil)
		assert.NotEqual(t, dbGroup.UpdatedAt, nil)
		assert.Equal(t, dbGroup.DeletedAt, nil)
		if dbGroup.ID <= 0 {
			t.Error("Id is not correctly assigned in the returned db object, got id: ", dbGroup.ID)
		}
	})

	t.Run("successfully retrieves a group by id", func(t *testing.T) {
		testHelper.TruncateTables([]string{"`group`"})
		groupName := "Group 1"
		dbGroup, _ := store.Group.CreateGroup(testHelper.Ctx, newGroup(groupName))
		retrievedGroup, err := store.Group.GetGroup(testHelper.Ctx, dbGroup.ID)
		assert.Equal(t, err, nil)
		assert.Equal(t, retrievedGroup.ID, dbGroup.ID)
		assert.Equal(t, *retrievedGroup.Name, *dbGroup.Name)
	})

	t.Run("successfully updates a group", func(t *testing.T) {
		testHelper.TruncateTables([]string{"`group`"})
		groupName := "Group 1"
		dbGroup, _ := store.Group.CreateGroup(testHelper.Ctx, newGroup(groupName))
		assert.Equal(t, *dbGroup.Name, groupName)
		newGroupName := "Group2"
		dbGroup.Name = &newGroupName
		dbGroup, err := store.Group.UpdateGroup(testHelper.Ctx, dbGroup)
		assert.Equal(t, err, nil)
		assert.Equal(t, *dbGroup.Name, newGroupName)
	})

	t.Run("Get number of profiles per group that were added before a specific time", func(t *testing.T) {
		testHelper.TruncateTables([]string{"group_profile", "profile", "`group`"})
		p1, _ := store.Profile.CreateProfile(testHelper.Ctx, newProfile("first", "last", "123"))
		p2, _ := store.Profile.CreateProfile(testHelper.Ctx, newProfile("first2", "last2", "12345"))
		g, _ := store.Group.CreateGroup(testHelper.Ctx, newGroup("group-name"))

		e := store.Profile.AddProfileToGroup(testHelper.Ctx, p1.ID, g.ID)
		assert.Equal(t, e, nil)
		e = store.Profile.AddProfileToGroup(testHelper.Ctx, p2.ID, g.ID)
		assert.Equal(t, e, nil)

		time.Sleep(time.Second * 1)

		result := store.Group.CountNumberOfProfilesToProcess(testHelper.Ctx, g.ID, time.Now().UTC())
		assert.Equal(t, result, 2)

	})
}
