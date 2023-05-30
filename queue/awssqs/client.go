package client

import (
	"customer_engagement/queue"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type sqsClient struct {
	// Referneces to AWS SQS client package. Used for direct communication with AWS SQS.
	client sqsiface.SQSAPI
}

// Constructs a new AWS SQS client.
func NewSqs(client sqsiface.SQSAPI) queue.IQueueClient {
	return sqsClient{
		client: client,
	}
}

// Sends a message request to the queue.
func (s sqsClient) Send(req *queue.SendRequest) (string, error) {
	attrs := make(map[string]*sqs.MessageAttributeValue, len(req.Attributes))
	for _, attr := range req.Attributes {
		attrs[attr.Key] = &sqs.MessageAttributeValue{
			StringValue: aws.String(attr.Value),
			DataType:    aws.String(attr.Type),
		}
	}

	request := &sqs.SendMessageInput{
		QueueUrl:          aws.String(req.QueueUrl),
		MessageBody:       aws.String(req.Body),
		MessageAttributes: attrs,
	}
	res, err := s.client.SendMessage(request)
	if err != nil {
		return "", fmt.Errorf("error sending message to the Queue: %w", err)
	}

	return *res.MessageId, nil
}

// Fetches a message from the requested queue.
func (s sqsClient) Receieve(queueUrl string) (*queue.Message, error) {
	res, err := s.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueUrl),
		MaxNumberOfMessages:   aws.Int64(1),
		WaitTimeSeconds:       aws.Int64(20),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	if err != nil {
		return nil, fmt.Errorf("error receiving message from the Queue: %w", err)

	}

	if len(res.Messages) == 0 {
		return nil, nil
	}

	return &queue.Message{
		ID:             *res.Messages[0].MessageId,
		ReceiptHandler: *res.Messages[0].ReceiptHandle,
		Body:           *res.Messages[0].Body,
		Attributes:     mapAttributes(res.Messages[0]),
	}, nil
}

// Deletes a message from the queue.
func (s sqsClient) Delete(queueUrl string, rcId string) error {
	deleteRequest := sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: &rcId,
	}
	_, err := s.client.DeleteMessage(&deleteRequest)
	if err != nil {
		return fmt.Errorf("error deleting message from the queue: %w", err)
	}
	return nil
}

func mapAttributes(m *sqs.Message) map[string]queue.Attribute {
	attrs := make(map[string]queue.Attribute)
	for k, v := range m.MessageAttributes {
		attrs[k] = queue.Attribute{
			Key:   k,
			Type:  *v.DataType,
			Value: *v.StringValue,
		}
	}

	return attrs

}
