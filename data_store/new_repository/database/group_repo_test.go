package dbRepository

import (
	testH "customer_engagement/test_helper"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

var (
	groupRepo IGroupRepository = NewGroupRepo(testH.DB)
)

func TestGroupRepository(t *testing.T) {
	testH.TruncateTables([]string{"`group`"})
	t.Run("successfully inserts a group", func(t *testing.T) {
		groupName := "Group 1"
		dbGroup, err := groupRepo.CreateGroup(testH.Ctx, newGroup(groupName))
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
		testH.TruncateTables([]string{"`group`"})
		groupName := "Group 1"
		dbGroup, _ := groupRepo.CreateGroup(testH.Ctx, newGroup(groupName))
		retrievedGroup, err := groupRepo.GetGroup(testH.Ctx, dbGroup.ID)
		assert.Equal(t, err, nil)
		assert.Equal(t, retrievedGroup.ID, dbGroup.ID)
		assert.Equal(t, *retrievedGroup.Name, *dbGroup.Name)
	})

	t.Run("successfully updates a group", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		groupName := "Group 1"
		dbGroup, _ := groupRepo.CreateGroup(testH.Ctx, newGroup(groupName))
		assert.Equal(t, *dbGroup.Name, groupName)
		newGroupName := "Group2"
		dbGroup.Name = &newGroupName
		dbGroup, err := groupRepo.UpdateGroup(testH.Ctx, dbGroup)
		assert.Equal(t, err, nil)
		assert.Equal(t, *dbGroup.Name, newGroupName)
	})
}

func newGroup(name string) *GroupStore { return &GroupStore{Name: &name} }
