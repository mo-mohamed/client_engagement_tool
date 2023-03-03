package queue

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type AwsQueueClient struct{}

type SendRequest struct {
	QueueURL   string
	Body       string
	Attributes []Attribute
}

type Attribute struct {
	Key   string
	Value string
	Type  string
}

func (*AwsQueueClient) Enqueue(request SendRequest) error {
	sqsClient := newSQS(os.Getenv("AWS_SQS_REGION"), os.Getenv("AWS_SQS_ENDPOINT"))
	_, err := sendMessage(sqsClient, request)
	return err
}

func newSQS(region, endpoint string) sqsiface.SQSAPI {
	cfg := aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}

	sess := session.Must(session.NewSession(&cfg))
	return sqs.New(sess)
}

func sendMessage(sqsClient sqsiface.SQSAPI, request SendRequest) (*sqs.SendMessageOutput, error) {

	attrs := make(map[string]*sqs.MessageAttributeValue, len(request.Attributes))
	for _, attr := range request.Attributes {
		attrs[attr.Key] = &sqs.MessageAttributeValue{
			StringValue: aws.String(attr.Value),
			DataType:    aws.String(attr.Type),
		}
	}

	sqsMessage := &sqs.SendMessageInput{
		QueueUrl:          aws.String(request.QueueURL),
		MessageBody:       aws.String(request.Body),
		MessageAttributes: attrs,
	}

	output, err := sqsClient.SendMessage(sqsMessage)
	if err != nil {
		return nil, fmt.Errorf("could not send message to queue %v: %v", request.QueueURL, err)
	}

	return output, nil
}
