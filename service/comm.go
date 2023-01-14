package service

import (
	communication "customer_engagement/comm"
	communicationClients "customer_engagement/comm/clients/sms"
	"fmt"
	"strconv"
)

type CommService struct{}

func (*CommService) Dispatch(groupId int) {
	go func() {
		pool := communication.NewComJobPool(2)
		var jobs []communication.ICommunication
		for i := 0; i < 10; i++ {
			j := communicationClients.NewSms("+12:"+strconv.Itoa(i), "hi there!")
			jobs = append(jobs, j)
		}

		go pool.AddBatch(jobs)
		pool.Run()

		for res := range pool.Results() {
			fmt.Println("done", res)
		}
	}()
}
