package consumers

import group "customer_engagement/consumers/group"

func NewGroupSMSConsumer() *group.GroupSmsConsumer {
	return group.NewGroupSMSConsumer()
}
