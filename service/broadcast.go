package service

import (
	"context"
	"customer_engagement/queue"
	interfaces "customer_engagement/service/interfaces"
	storeLayer "customer_engagement/store"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type BroadCastService struct {
	queueClient queue.IQueueClient
	store       *storeLayer.Store
}

func NewBroadcastService(store *storeLayer.Store, queueClient queue.IQueueClient) interfaces.IBroadcastService {
	return &BroadCastService{queueClient: queueClient, store: store}
}

func (b *BroadCastService) EnqueueBroadcastSimpleSmsToGroup(ctx context.Context, message string, groupId int) (*string, error) {
	exists, _ := b.store.Group.Exists(ctx, groupId)
	if !exists {
		return nil, fmt.Errorf("group with ID %v is not found ", groupId)
	}

	attributes := make([]queue.Attribute, 0)
	attributes = append(attributes, queue.Attribute{
		Key:   "GroupId",
		Value: strconv.Itoa(groupId),
		Type:  "String",
	})
	attributes = append(attributes, queue.Attribute{
		Key:   "DateEnqueued",
		Value: time.Now().UTC().String(),
		Type:  "String",
	})
	attributes = append(attributes, queue.Attribute{
		Key:   "InternalID",
		Value: uuid.NewString(),
		Type:  "String",
	})

	messageRequest := &queue.SendRequest{
		QueueUrl:   os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Body:       message,
		Attributes: attributes,
	}

	messageId, err := b.queueClient.Send(messageRequest)
	if err != nil {
		fmt.Println("Error")
		return nil, fmt.Errorf("error occured while enqueing")
	}

	return messageId, err
}
