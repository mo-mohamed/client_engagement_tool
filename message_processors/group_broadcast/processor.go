package groupBroadcastProcessor

import (
	"customer_engagement/queue"
	"fmt"
)

type GroupMessageProcessor struct{}

func (g *GroupMessageProcessor) Process(message *queue.Message) error {

	// Write down group logic here
	groupid := message.Attributes["GroupId"]
	fmt.Println("here is the group", groupid)
	return nil
}
