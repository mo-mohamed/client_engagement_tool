package interfaces

import "context"

type IBroadcastService interface {
	EnqueueBroadcastSimpleSmsToGroup(ctx context.Context, message string, groupId int) (*string, error)
}
