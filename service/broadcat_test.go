package service

import (
	mock_queue "customer_engagement/mocks/queue"
	interfaces "customer_engagement/service/interfaces"
	domains "customer_engagement/service/models"
	storeLayer "customer_engagement/store"
	storeRepo "customer_engagement/store/repository"
	testH "customer_engagement/test_helper"
	"testing"

	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestBroadcastService(t *testing.T) {

	ctrl := gomock.NewController(t)
	client := mock_queue.NewMockIQueueClient(ctrl)
	defer ctrl.Finish()

	var storedb *storeLayer.Store = &storeLayer.Store{
		Profile: storeRepo.NewProfileRepo(testH.DB),
		Group:   storeRepo.NewGroupRepo(testH.DB),
	}
	var broadcastService interfaces.IBroadcastService = NewBroadcastService(storedb, client)
	var groupService interfaces.IGroupService = NewGroupService(storedb)

	t.Run("successfully broadcasts a messaage", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		group_name := "group-name"
		g, err := groupService.Create(testH.Ctx, &domains.Group{Name: &group_name})
		assert.Equal(t, err, nil)

		// Fake Request
		messageId := "id_queue_broadcast"
		client.EXPECT().Send(gomock.Any()).Return(&messageId, nil)

		res, err := broadcastService.EnqueueBroadcastSimpleSmsToGroup(testH.Ctx, "message body", g.ID)
		assert.Equal(t, *res, messageId)
		assert.Equal(t, err, nil)
	})

	t.Run("returns an error if a group was not found", func(t *testing.T) {
		testH.TruncateTables([]string{"`group`"})
		groupId := 5

		res, err := broadcastService.EnqueueBroadcastSimpleSmsToGroup(testH.Ctx, "message body", groupId)
		assert.Equal(t, res, nil)
		assert.NotEqual(t, err, nil)
	})
}
