package comm

type IJobPool interface {
	Run()
	Results() <-chan Response
	AddBatch([]ICommunication)
	Stop()
}

type ICommunication interface {
	Send() error
}
