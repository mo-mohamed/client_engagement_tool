package mock_interfaces

import (
	"customer_engagement/queue"
)

type MessageProcessorMock struct {
	Processed bool
}

func (p *MessageProcessorMock) Process(*queue.Message) error {
	p.Processed = true
	return nil
}
