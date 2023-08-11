package interfaces

type IBroadcastService interface {
	EnqueueBroadcastSimpleSmsToGroup(message string, groupId int) (string, error)
}
