package sqsClient

import (
	"customer_engagement/queue"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsClient struct {
	// Referneces to AWS SQS client package. Used for direct communication with AWS SQS.
	client *sqs.SQS
}

// Constructs a new AWS SQS client.
func NewSqs(session *session.Session) queue.IQueueClient {
	return sqsClient{
		client: sqs.New(session),
	}
}

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

	attrs := make([]queue.Attribute, 0)
	for k, v := range res.Messages[0].MessageAttributes {
		attrs = append(attrs, queue.Attribute{
			Key:   k,
			Value: *v.StringValue,
			Type:  *v.DataType,
		})
	}

	return &queue.Message{
		ID:             *res.Messages[0].MessageId,
		ReceiptHandler: *res.Messages[0].ReceiptHandle,
		Body:           *res.Messages[0].Body,
		Attributes:     attrs,
	}, nil
}

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
