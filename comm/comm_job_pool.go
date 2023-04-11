package comm

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type ComJobPool struct {
	concurrency   int
	jobsStream    chan ICommunication
	resultsStream chan Response
	done          chan interface{}
	ctx           context.Context
	ctxCancelFn   context.CancelFunc
}

func NewComJobPool(concurrency int) *ComJobPool {
	ctx, ctxCancelFn := context.WithCancel(context.TODO())
	return &ComJobPool{
		concurrency:   concurrency,
		jobsStream:    make(chan ICommunication, concurrency),
		resultsStream: make(chan Response, concurrency),
		done:          make(chan interface{}),
		ctx:           ctx,
		ctxCancelFn:   ctxCancelFn,
	}
}

func (p *ComJobPool) Stop() {
	p.ctxCancelFn()
}

func (p *ComJobPool) AddBatch(jobs []ICommunication) {
	for _, job := range jobs {
		p.jobsStream <- job
	}
	//done writing to the jobs channel, close it
	close(p.jobsStream)
}

func (p *ComJobPool) Results() <-chan Response {
	return p.resultsStream
}

func (p *ComJobPool) Run() {
	var wg sync.WaitGroup
	for i := 0; i < p.concurrency; i++ {
		wg.Add(1)
		go p.doWork(&wg)

	}

	wg.Wait()
	// work completed. close all channels
	close(p.resultsStream)
	close(p.done)

}

func (p *ComJobPool) doWork(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-p.jobsStream:
			if !ok {
				// reading from a closed channel, all jobs are taken
				return
			}
			err := job.Process()
			if err != nil {
				fmt.Println("error")
				p.resultsStream <- Response{Status: 500, Error: errors.New("err")}
			} else {
				fmt.Println("success")
				p.resultsStream <- Response{Status: 200, Error: nil}
			}
		case <-p.ctx.Done():
			// p.ResultsStream <- models.Response{Status: 100, Error: errors.New("cancelled")}
			return
		}
	}

}
