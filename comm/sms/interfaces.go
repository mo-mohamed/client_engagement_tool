package sms

type ISms interface {
	Send(Request) Response
}

type ISmsProcessor interface {
	Run()
	Results() <-chan Response
	AddBatch([]Request)
	Stop()
}
