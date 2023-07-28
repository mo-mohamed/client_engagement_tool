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

var (
	store_db *storeLayer.Store = &storeLayer.Store{
		Profile: storeRepo.NewProfileRepo(testH.DB),
		Group:   storeRepo.NewGroupRepo(testH.DB),
	}
	groupService interfaces.IGroupService = NewGroupService(store_db)
)

func TestGroupService(t *testing.T) {
	t.Run("successfully create a group", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		_, err := groupService.Create(testH.Ctx, newGroup("group 1"))
		assert.Equal(t, err, nil)
	})

	t.Run("successfully gets a group via id", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		g, _ := groupService.Create(testH.Ctx, newGroup("group 2"))
		g, err := groupService.Get(testH.Ctx, g.ID)
		assert.Equal(t, err, nil)
		assert.Equal(t, *g.Name, "group 2")
	})

	t.Run("successfully updates a group", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		g, _ := groupService.Create(testH.Ctx, newGroup("group name"))
		newName := "group-updated"
		g.Name = &newName
		g, err := groupService.Update(testH.Ctx, g)
		assert.Equal(t, err, nil)
		assert.Equal(t, *g.Name, newName)
	})

	t.Run("checks existence of a group", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		g, _ := groupService.Create(testH.Ctx, newGroup("group"))
		exists, err := groupService.Exists(testH.Ctx, g.ID)
		assert.Equal(t, err, nil)
		assert.Equal(t, exists, true)
	})

	t.Run("checks existence of a group when a group doesn't exists", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		exists, err := groupService.Exists(testH.Ctx, 2)
		assert.Equal(t, err, nil)
		assert.Equal(t, exists, true)
	})
}

func newGroup(name string) *domains.Group { return &domains.Group{Name: &name} }
