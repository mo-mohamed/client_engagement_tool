package consumers

import (
	messageConsumer "customer_engagement/consumers/message"
	sqsclient "customer_engagement/queue/awssqs"
	"os"
	"strconv"

	interfaces "customer_engagement/consumers/interfaces"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	gbc "customer_engagement/message_processors/group_broadcast"
)

func NewGroupQueueConsumer() interfaces.IConsumer {

	sqsSession := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_SQS_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_SQS_ENDPOINT")),
	}))

	queueConfig := messageConsumer.QueueConfig{
		Url:    os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Client: sqsclient.NewSqs(sqs.New(sqsSession)),
	}

	profilesBatchSize, _ := strconv.Atoi(os.Getenv("PROFILES_BATCH_SIZE"))
	processor := gbc.NewGroupMessageProcessor(profilesBatchSize)

	numOfGroupWorkers, _ := strconv.Atoi(os.Getenv("NUM_GROUP_MESSAGE_QUEUE_CONSUMER"))
	return messageConsumer.NewMessageConsumer(numOfGroupWorkers, queueConfig, processor)
}
