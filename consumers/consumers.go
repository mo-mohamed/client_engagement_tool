package consumers

import (
	messageConsumer "customer_engagement/consumers/message"
	sqsclient "customer_engagement/queue/awssqs"
	"os"
	"strconv"

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
	client := sqsclient.NewSqs(sess)
	queueConfig := messageConsumer.QueueConfig{
		Url:    os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Client: &client,
	}

	profilesBatchSize, _ := strconv.Atoi(os.Getenv("PROFILES_BATCH_SIZE"))
	processor := gbc.NewGroupMessageProcessor(profilesBatchSize)

	numOfGroupWorkers, _ := strconv.Atoi(os.Getenv("NUM_GROUP_MESSAGE_QUEUE_CONSUMER"))
	return messageConsumer.NewMessageConsumer(numOfGroupWorkers, queueConfig, processor)
}
