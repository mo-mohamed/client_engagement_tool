package client

import (
	mock_sqsiface "customer_engagement/mocks/aws_sqs"
	"customer_engagement/queue"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestSuccessEnqueueOfAMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock_sqsiface.NewMockSQSAPI(ctrl)
	defer ctrl.Finish()

	sqsClient := NewSqs(client)

	send_request := queue.SendRequest{
		QueueUrl: "path-for-queue",
		Body:     "message-body",
	}
	messageId := "id_queue_123"

	request := &sqs.SendMessageInput{
		QueueUrl:          aws.String("path-for-queue"),
		MessageBody:       aws.String("message-body"),
		MessageAttributes: make(map[string]*sqs.MessageAttributeValue, 0),
	}

	client.EXPECT().SendMessage(request).Return(&sqs.SendMessageOutput{MessageId: &messageId}, nil)
	result, err := sqsClient.Send(&send_request)
	assert.Equal(t, result, "id_queue_123")
	assert.Equal(t, err, nil)
}

func TestCapturesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_sqsiface.NewMockSQSAPI(ctrl)
	defer ctrl.Finish()

	sqsClient := NewSqs(mockClient)

	send_request := queue.SendRequest{
		QueueUrl: "path-for-queue",
		Body:     "message-body",
	}

	request := &sqs.SendMessageInput{
		QueueUrl:          aws.String("path-for-queue"),
		MessageBody:       aws.String("message-body"),
		MessageAttributes: make(map[string]*sqs.MessageAttributeValue, 0),
	}

	mockClient.EXPECT().SendMessage(request).Return(nil, errors.New("failed to communicate with sqs service"))
	result, err := sqsClient.Send(&send_request)
	assert.Equal(t, result, "")
	assert.Equal(t, err.Error(), "error sending message to the Queue: failed to communicate with sqs service")
}
