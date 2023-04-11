package queue

type IQueueClient interface {
	// Sends a request to the message queue.
	// Accepts a `SendRequest` pointer.
	Send(req *SendRequest) (string, error)

	// Receieves a single message from the queue.
	// accepts a queue url.
	Receieve(queueUrl string) (*Message, error)

	// Deletes a message from the queue.
	// accepts a message handler.
	Delete(queueUrl string, rcId string) error
}
