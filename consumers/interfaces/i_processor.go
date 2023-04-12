package interfaces

import (
	"customer_engagement/queue"
)

/*
This interface is require dto be implemented by the processor, which is passed to the consumer
to do the actual processing from the messages received from the queue.
*/
type IProcessor interface {

	// Process a messahe received from the consumer's queue client
	Process(*queue.Message) error
}
