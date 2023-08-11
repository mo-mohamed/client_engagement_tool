package service

import (
	"customer_engagement/queue"
	interfaces "customer_engagement/service/interfaces"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type BroadCastService struct {
	queueClient queue.IQueueClient
}

func NewBroadcastService(queueClient queue.IQueueClient) interfaces.IBroadcastService {
	return &BroadCastService{queueClient: queueClient}
}

func (b *BroadCastService) EnqueueBroadcastSimpleSmsToGroup(message string, groupId int) (string, error) {
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

	return b.queueClient.Send(messageRequest)
}
