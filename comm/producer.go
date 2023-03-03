package comm

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type GroupProducer struct{}

func (*GroupProducer) EnqueueGroupMessage(group_id, body string) {
	sqsClient := newSQS(os.Getenv("AWS_SQS_REGION"), os.Getenv("AWS_SQS_ENDPOINT"))
	_, err := sendMessage(sqsClient, "helloooooooo", os.Getenv("AWS_SQS_ENDPOINT")+"/"+os.Getenv("AWS_SQS_SMS_GROUP_NAME"))
	if err != nil {
		fmt.Println(err)
	}
}

func newSQS(region, endpoint string) sqsiface.SQSAPI {
	cfg := aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}

	sess := session.Must(session.NewSession(&cfg))
	return sqs.New(sess)
}

func sendMessage(sqsClient sqsiface.SQSAPI, msg, queueURL string) (*sqs.SendMessageOutput, error) {
	sqsMessage := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(msg),
	}

	output, err := sqsClient.SendMessage(sqsMessage)
	if err != nil {
		return nil, fmt.Errorf("could not send message to queue %v: %v", queueURL, err)
	}

	return output, nil
}
