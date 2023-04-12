package messageConsumer

import (
	"customer_engagement/queue"
	"runtime"
	"testing"
	"time"
)

type queueClientMock struct{}

func (q queueClientMock) Send(req *queue.SendRequest) (string, error) {
	return "", nil
}

func (q queueClientMock) Receieve(queueUrl string) (*queue.Message, error) {
	return &queue.Message{}, nil
}

func (q queueClientMock) Delete(queueUrl string, rcId string) error {
	return nil
}

type processorMock struct {
	processed bool
}

func (p *processorMock) Process(*queue.Message) error {
	p.processed = true
	return nil
}

func TestConsumer(t *testing.T) {
	var clientMock queue.IQueueClient = queueClientMock{}

	queueConfig := QueueConfig{
		Url:    "TestURL",
		Client: &clientMock,
	}

	processorMock := processorMock{processed: false}
	c := NewMessageConsumer(1, queueConfig, &processorMock)
	go c.Run()
	go func() {
		time.Sleep(time.Millisecond * 10)
		c.Stop()
	}()
	time.Sleep(time.Millisecond * 10)
	if !processorMock.processed {
		t.Error("Processor is not completed.")
	}
}

func TestNumberOfConsurrency(t *testing.T) {
	current := runtime.NumGoroutine()
	concurrency := 3
	var clientMock queue.IQueueClient = queueClientMock{}

	queueConfig := QueueConfig{
		Url:    "TestURL",
		Client: &clientMock,
	}
	processorMock := processorMock{}
	c := NewMessageConsumer(concurrency, queueConfig, &processorMock)
	go c.Run()

	time.Sleep(time.Millisecond * 10)
	if runtime.NumGoroutine() != current+concurrency+1 { // +1 for the go c.Run()
		t.Error("Number of go routines is not as expected")
	}

}
