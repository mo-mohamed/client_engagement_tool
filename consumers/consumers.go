package consumers

import (
	gconsumer "customer_engagement/consumers/message"
	sqsclient "customer_engagement/queue/awssqs"
	"os"

	interfaces "customer_engagement/consumers/interfaces"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	gbc "customer_engagement/message_processors/group_broadcast"
)

func NewGroupQueueConsumer() interfaces.IConsumer {
	cfg := aws.Config{
		Region:   aws.String(os.Getenv("AWS_SQS_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_SQS_ENDPOINT")),
	}

	sess := session.Must(session.NewSession(&cfg))
	s := sqsclient.NewSqs(sess)
	queueConfig := gconsumer.QueueConfig{
		Url:    os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Client: &s,
	}

	pro := gbc.GroupMessageProcessor{}

	return gconsumer.NewMessageConsumer(2, queueConfig, &pro)
}
