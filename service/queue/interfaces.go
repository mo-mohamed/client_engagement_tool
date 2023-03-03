package queue

type IGroupQueue interface {
	Enqueue(request SendRequest) error
}
