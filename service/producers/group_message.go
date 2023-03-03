package producers

import (
	"customer_engagement/service/queue"
	"os"
	"strconv"
)

type GroupMessageProducer struct{}

func (*GroupMessageProducer) EnqueueGroupBroadcast(groupId int, message string) error {
	attrs := make([]queue.Attribute, 0)
	attrs = append(attrs, queue.Attribute{
		Key:   "groupId",
		Value: strconv.Itoa(groupId),
		Type:  "String",
	})

	req := queue.SendRequest{
		QueueURL:   os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Body:       message,
		Attributes: attrs,
	}

	return queue.GetClient().Enqueue(req)
}
