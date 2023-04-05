package sqsClient

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type AwsQueue struct {
	client sqsiface.SQSAPI
}

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

func New(region, endpoint string) AwsQueue {
	cfg := aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}
	sess := session.Must(session.NewSession(&cfg))
	return AwsQueue{
		client: sqs.New(sess),
	}
}

func (awsQueue *AwsQueue) Enqueue(request SendRequest) error {
	_, err := sendMessage(awsQueue.client, request)
	return err
}

func sendMessage(sqsClient sqsiface.SQSAPI, request SendRequest) (*sqs.SendMessageOutput, error) {

	attrs := make(map[string]*sqs.MessageAttributeValue, len(request.Attributes))
	for _, attr := range request.Attributes {
		attrs[attr.Key] = &sqs.MessageAttributeValue{
			StringValue: aws.String(attr.Value),
			DataType:    aws.String(attr.Type),
		}
	}

	fmt.Println(attrs)

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

func (awsQueue *AwsQueue) ReceiveMessage(r *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return awsQueue.client.ReceiveMessage(r)
}

func (awsQueue *AwsQueue) DeleteMessage(d *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return awsQueue.client.DeleteMessage(d)
}
