package service

import (
	interfaces "customer_engagement/service/interfaces"
)

type Service struct {
	Group   interfaces.IGroupService
	Profile interfaces.IProfileService
}
