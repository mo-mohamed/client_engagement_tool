package producers

import (
	"customer_engagement/service/queue"
	"fmt"
	"os"
	"strconv"
	"time"
)

type GroupMessageProducer struct{}

func (*GroupMessageProducer) EnqueueGroupBroadcast(groupId int, message string) error {
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

	req := queue.SendRequest{
		QueueURL:   os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Body:       message,
		Attributes: attributes,
	}

	fmt.Println("send request:", req)
	err := queue.GetClient().Enqueue(req)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
