package producers

import (
	sqsClient "customer_engagement/clients/sqs"
	"fmt"
	"os"
	"strconv"
	"time"
)

type GroupMessageProducer struct{}

func (*GroupMessageProducer) EnqueueGroupBroadcast(groupId int, message string) error {
	attributes := make([]sqsClient.Attribute, 0)
	attributes = append(attributes, sqsClient.Attribute{
		Key:   "GroupId",
		Value: strconv.Itoa(groupId),
		Type:  "String",
	})
	attributes = append(attributes, sqsClient.Attribute{
		Key:   "DateEnqueued",
		Value: time.Now().UTC().String(),
		Type:  "String",
	})

	req := sqsClient.SendRequest{
		QueueURL:   os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Body:       message,
		Attributes: attributes,
	}

	client := sqsClient.New(os.Getenv("AWS_SQS_REGION"), os.Getenv("AWS_SQS_ENDPOINT"))
	err := client.Enqueue(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
