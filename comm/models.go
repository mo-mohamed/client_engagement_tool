package sms

type Request struct {
	Destination string
	Body        string
	// ExecFn      func() error
}

type Response struct {
	Error  error
	Status int
}
