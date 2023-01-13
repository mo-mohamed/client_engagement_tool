package sms

import (
	interfaces "customer_engagement/comm/clients"
)

type ISms interface {
	Send(Request) Response
}

type ISmsProcessor interface {
	Run()
	Results() <-chan Response
	AddBatch([]interfaces.ICommunication)
	Stop()
}
