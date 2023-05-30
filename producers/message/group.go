package producer

import (
	"customer_engagement/queue"
)

type GroupProducer struct {
	client queue.IQueueClient
}

/*
Returns a new instance of the group processor.
Accepts an IQueueClient as a paramater
*/
func NewGroupProducer(c queue.IQueueClient) *GroupProducer {
	return &GroupProducer{
		client: c,
	}
}

/*
Calls the underlying broker client to send a message to its queue.
*/
func (g *GroupProducer) Produce(request *queue.SendRequest) (string, error) {
	return g.client.Send(request)
}
