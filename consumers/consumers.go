package consumers

import group "customer_engagement/consumers/group"

func NewGroupQueueConsumer() *group.GroupQueueConsumer {
	return group.NewGroupQueueConsumer()
}
