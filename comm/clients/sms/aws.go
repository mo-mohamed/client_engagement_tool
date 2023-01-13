package aws

import (
	"fmt"
)

type sms struct {
	MDN  string
	Body string
}

func NewSms(mdn string, body string) *sms {
	return &sms{
		MDN:  mdn,
		Body: body,
	}
}

func (s *sms) Send() error {

	// simulate work
	fmt.Println("Done sending a message to", s.MDN, "and body is", s.Body)
	return nil
}
