package sms

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type SmsProcessor struct {
	Concurrency   int
	JobsStream    chan Request
	ResultsStream chan Response
	Done          chan interface{}
	ctx           context.Context
	ctxCancelFn   context.CancelFunc
}

func NewSmsProcessor(concurrency int) ISmsProcessor {
	ctx, ctxCancelFn := context.WithCancel(context.TODO())
	return &SmsProcessor{
		Concurrency:   concurrency,
		JobsStream:    make(chan Request, concurrency),
		ResultsStream: make(chan Response, concurrency),
		Done:          make(chan interface{}),
		ctx:           ctx,
		ctxCancelFn:   ctxCancelFn,
	}
}

func (p *SmsProcessor) Stop() {
	p.ctxCancelFn()
}

func (p *SmsProcessor) AddBatch(jobs []Request) {
	for i, job := range jobs {
		p.JobsStream <- job
		fmt.Println("Added job", i)
	}

	//done writing to the jobs channel, close it out
	fmt.Println("closing channel")
	close(p.JobsStream)
	fmt.Println("done closing channel")
}

func (p *SmsProcessor) Results() <-chan Response {
	return p.ResultsStream
}

func (p *SmsProcessor) Run() {
	var wg sync.WaitGroup
	for i := 0; i < p.Concurrency; i++ {
		wg.Add(1)
		go p.doWork(&wg)
	}

	wg.Wait()
	fmt.Println("done waiting")
	// work completed. close all channels
	fmt.Println(p.ResultsStream)
	fmt.Println(p.Done)
	close(p.ResultsStream)
	close(p.Done)

}

func (p *SmsProcessor) doWork(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Spawning")
	defer func() {
		fmt.Println("Exited")
	}()
	for {
		select {
		case job, ok := <-p.JobsStream:
			if !ok {
				// reading from an empty channel
				return
			}
			fmt.Println("got job", job)
			job.ExecFn()
			// simulate work (should be sending SMS)
			time.Sleep(time.Second * 1)

			// do the actual work and send the results back to the channel
			p.ResultsStream <- Response{Status: 500, Error: errors.New("err")}
		case <-p.ctx.Done():
			fmt.Println("received done request")
			// p.ResultsStream <- models.Response{Status: 100, Error: errors.New("cancelled")}
			return
		}
	}

}
