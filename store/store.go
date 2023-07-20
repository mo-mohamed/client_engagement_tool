package store

import (
	iStore "customer_engagement/store/interfaces"
)

type Store struct {
	Profile iStore.IProfileRepo
	Group   iStore.IGroupRepository
}
