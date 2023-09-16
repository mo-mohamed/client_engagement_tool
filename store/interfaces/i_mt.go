package interfaces

import (
	"context"
	"customer_engagement/store/models"
)

type IMobileTerminated interface {
	CreateMT(ctx context.Context, mt *models.MT) (*models.MT, error)
	IsProcessed(ctx context.Context, mt *models.MT) bool
}
