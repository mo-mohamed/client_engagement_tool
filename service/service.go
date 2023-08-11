package service

import (
	service_i "customer_engagement/service/interfaces"
)

type Service struct {
	Group     service_i.IGroupService
	Profile   service_i.IProfileService
	Broadcast service_i.IBroadcastService
}
