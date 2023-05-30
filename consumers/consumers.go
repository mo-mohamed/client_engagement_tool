package consumers

import (
	group "customer_engagement/consumers/message"
	sqsClient "customer_engagement/queue/awssqs"
	"os"

	int "customer_engagement/consumers/interfaces"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	gbc "customer_engagement/message_processors/group_broadcast"
)

func NewGroupQueueConsumer() int.IConsumer {
	cfg := aws.Config{
		Region:   aws.String(os.Getenv("AWS_SQS_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_SQS_ENDPOINT")),
	}

	sess := session.Must(session.NewSession(&cfg))
	s := sqsClient.NewSqs(sess)
	v := group.QueueConfig{
		Url:    os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Client: &s,
	}

	pro := gbc.GroupMessageProcessor{}

	return group.NewMessageConsumer(2, v, &pro)
}
