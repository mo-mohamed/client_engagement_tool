package interfaces

import (
	"context"
	domains "customer_engagement/service/models"
)

type IMobileTerminated interface {
	create(ctx context.Context, mt domains.MT) (*domains.MT, error)
	IsProcessed(ctx context.Context, mt domains.MT) bool
}
