package sms

import (
	"customer_engagement/comm/sms"
	"fmt"
)

type AwsSMS struct{}

func (s AwsSMS) Send(request sms.Request) sms.Response {

	// simulate work
	fmt.Println("Done sending a message to", request.Destination, request.Body)
	return sms.Response{
		Error:  nil,
		Status: 200,
	}
}
