package client

import (
	"customer_engagement/queue"
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"gopkg.in/go-playground/assert.v1"
)

type mockSQSClient struct {
	sqsiface.SQSAPI
}

func (m *mockSQSClient) SendMessage(in *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	messageId := "id_queue_123"
	return &sqs.SendMessageOutput{MessageId: &messageId}, nil
}
func (m *mockSQSClient) ReceiveMessage(in *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return &sqs.ReceiveMessageOutput{}, nil
}

func ReturnsMessageIdAfterSendingAMessage(t *testing.T) {
	sqsClient := NewSqs(&mockSQSClient{})
	attributes := make([]queue.Attribute, 0)
	attributes = append(attributes, queue.Attribute{
		Key:   "GroupId",
		Value: "123",
		Type:  "String",
	})
	send_request := queue.SendRequest{
		QueueUrl:   "path-for-queue",
		Body:       "message-body",
		Attributes: attributes,
	}
	result, err := sqsClient.Send(&send_request)
	assert.Equal(t, result, "id_queue_123")
	assert.Equal(t, err, nil)
	// queueURL := "https://queue.amazonaws.com/80398EXAMPLE/MyQueue"
	// q.SendMessage(&sqs.SendMessageInput{
	// 	MessageBody: aws.String("Hello, World!"),
	// 	QueueUrl:    &queueURL,
	// })
	// message, _ := q.ReceiveMessage(&sqs.ReceiveMessageInput{
	// 	QueueUrl: &queueURL,
	// })
	// assert.Equal(t, *message.Messages[0].Body, "Hello, World!")
}
