package messageConsumer

import (
	mockProcessor "customer_engagement/mocks/consumer_processor"
	mockQueue "customer_engagement/mocks/queue"
	"customer_engagement/queue"
	"runtime"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestMessageConsumer(t *testing.T) {
	t.Run("successfully process a request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		queuetMock := mockQueue.NewMockIQueueClient(ctrl)
		defer ctrl.Finish()

		queueConfig := QueueConfig{
			Url:    "queue-url",
			Client: queuetMock,
		}

		ret := &queue.Message{
			ID:             "id",
			ReceiptHandler: "handler-id",
			Body:           "body",
			Attributes: map[string]queue.Attribute{
				"GroupId": queue.Attribute{
					Key:   "GroupId",
					Value: "1",
					Type:  "string",
				},
			},
		}

		processorMock := &mockProcessor.MessageProcessorMock{Processed: false}
		consumer := NewMessageConsumer(1, queueConfig, processorMock)
		queuetMock.EXPECT().Receieve("queue-url").Return(ret, nil).AnyTimes()
		queuetMock.EXPECT().Delete("queue-url", "handler-id").Return(nil).AnyTimes()
		go consumer.Run()
		go func() {
			time.Sleep(time.Millisecond * 10)
			consumer.Stop()
		}()
		time.Sleep(time.Millisecond * 10)
		assert.Equal(t, processorMock.Processed, true)
	})

	t.Run("Spins up the correct number of go routines as workers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		queuetMock := mockQueue.NewMockIQueueClient(ctrl)
		defer ctrl.Finish()

		currentNumRoutines := runtime.NumGoroutine()
		numOfWorkers := 3

		returnQueueRequest := &queue.Message{
			ID:             "id",
			ReceiptHandler: "handler-id",
			Body:           "body",
			Attributes: map[string]queue.Attribute{
				"GroupId": queue.Attribute{
					Key:   "GroupId",
					Value: "1",
					Type:  "string",
				},
			},
		}
		queuetMock.EXPECT().Receieve("queue-url").Return(returnQueueRequest, nil).AnyTimes()
		queuetMock.EXPECT().Delete("queue-url", "handler-id").Return(nil).AnyTimes()
		queueConfig := QueueConfig{
			Url:    "queue-url",
			Client: queuetMock,
		}
		processorMock := mockProcessor.MessageProcessorMock{}
		c := NewMessageConsumer(numOfWorkers, queueConfig, &processorMock)
		go c.Run()

		time.Sleep(time.Millisecond * 10)
		expectedNumberOfRoutines := currentNumRoutines + numOfWorkers + 1
		assert.Equal(t, runtime.NumGoroutine(), expectedNumberOfRoutines)
	})
}
