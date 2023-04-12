package producer

import (
	"customer_engagement/queue"
)

type GroupProducer struct {
	client queue.IQueueClient
}

func NewGroupProducer(c queue.IQueueClient) *GroupProducer {
	return &GroupProducer{
		client: c,
	}
}

func (g *GroupProducer) Produce(request *queue.SendRequest) (string, error) {
	return g.client.Send(request)
}
