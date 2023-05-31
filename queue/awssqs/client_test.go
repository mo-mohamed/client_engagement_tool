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

func TestSendMessage(t *testing.T) {
	t.Run("Success sending a message to the queue", func(t *testing.T) {
		t.Parallel()
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
	})

	t.Run("Capture error when sending a message", func(t *testing.T) {
		t.Parallel()
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
	})

}

func TestReceiveMessage(t *testing.T) {
	t.Run("Receives a message", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := mock_sqsiface.NewMockSQSAPI(ctrl)
		defer ctrl.Finish()
		sqsClient := NewSqs(mockClient)

		input := &sqs.ReceiveMessageInput{
			QueueUrl:              aws.String("queue-url"),
			MaxNumberOfMessages:   aws.Int64(1),
			WaitTimeSeconds:       aws.Int64(20),
			MessageAttributeNames: aws.StringSlice([]string{"All"}),
		}

		toReturn := &sqs.ReceiveMessageOutput{
			Messages: []*sqs.Message{{
				Attributes:    make(map[string]*string, 0),
				Body:          aws.String("body"),
				MessageId:     aws.String("message-id"),
				ReceiptHandle: aws.String("handler-id"),
			}},
		}
		mockClient.EXPECT().ReceiveMessage(input).Return(toReturn, nil)
		result, err := sqsClient.Receieve("queue-url")
		assert.Equal(t, result.ID, "message-id")
		assert.Equal(t, result.Body, "body")
		assert.Equal(t, result.ReceiptHandler, "handler-id")
		assert.Equal(t, err, nil)
	})

	t.Run("Captures error on receiving a message", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := mock_sqsiface.NewMockSQSAPI(ctrl)
		defer ctrl.Finish()
		sqsClient := NewSqs(mockClient)

		input := &sqs.ReceiveMessageInput{
			QueueUrl:              aws.String("queue-url"),
			MaxNumberOfMessages:   aws.Int64(1),
			WaitTimeSeconds:       aws.Int64(20),
			MessageAttributeNames: aws.StringSlice([]string{"All"}),
		}
		mockClient.EXPECT().ReceiveMessage(input).Return(nil, errors.New("can't communicate with AWS"))
		result, err := sqsClient.Receieve("queue-url")
		assert.Equal(t, result, nil)
		assert.Equal(t, err.Error(), "error receiving message from the Queue: can't communicate with AWS")
	})
}

func TestDeleteMessage(t *testing.T) {
	t.Run("Deletes a message", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := mock_sqsiface.NewMockSQSAPI(ctrl)
		defer ctrl.Finish()
		sqsClient := NewSqs(mockClient)

		deleteRequest := sqs.DeleteMessageInput{
			QueueUrl:      aws.String("queue-url"),
			ReceiptHandle: aws.String("handler-id"),
		}

		mockClient.EXPECT().DeleteMessage(&deleteRequest).Return(&sqs.DeleteMessageOutput{}, nil)

		err := sqsClient.Delete("queue-url", "handler-id")
		assert.Equal(t, err, nil)
	})

	t.Run("Captures error on delete messages failure", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockClient := mock_sqsiface.NewMockSQSAPI(ctrl)
		defer ctrl.Finish()
		sqsClient := NewSqs(mockClient)

		deleteRequest := sqs.DeleteMessageInput{
			QueueUrl:      aws.String("queue-url"),
			ReceiptHandle: aws.String("handler-id"),
		}

		mockClient.EXPECT().DeleteMessage(&deleteRequest).Return(&sqs.DeleteMessageOutput{}, errors.New("some error"))

		err := sqsClient.Delete("queue-url", "handler-id")
		assert.Equal(t, err.Error(), "error deleting message from the queue: some error")
	})
}
