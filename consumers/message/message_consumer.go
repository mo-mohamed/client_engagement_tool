package messageConsumer

import (
	"context"
	"fmt"
	"sync"

	"customer_engagement/consumers/interfaces"
	qu "customer_engagement/queue"
)

type QueueConfig struct {
	// Full url of the queue
	Url string

	// A client that can interact with the queue
	Client *qu.IQueueClient
}

type QueueConsumer struct {
	concurrency int
	done        chan interface{}
	ctx         context.Context
	ctxCancelFn context.CancelFunc
	client      qu.IQueueClient
	queueUrl    string
	processor   interfaces.IProcessor
}

/*
Creates a new message consumer.

Returns in instance of the interface `interfaces.IProcessor`
*/
func NewMessageConsumer(numOfWorkers int, config QueueConfig, processor interfaces.IProcessor) interfaces.IConsumer {
	ctx, ctxCancelFn := context.WithCancel(context.TODO())
	c := &QueueConsumer{
		concurrency: numOfWorkers,
		done:        make(chan interface{}),
		ctx:         ctx,
		ctxCancelFn: ctxCancelFn,
		client:      *config.Client,
		queueUrl:    config.Url,
		processor:   processor,
	}
	return c
}

func (c *QueueConsumer) Stop() {
	c.ctxCancelFn()
}

func (c *QueueConsumer) Run() {
	var wg sync.WaitGroup
	for i := 0; i < c.concurrency; i++ {
		wg.Add(1)
		go c.doWork(&wg, i)

	}
	wg.Wait()
	close(c.done)
}

func (c *QueueConsumer) doWork(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			res, err := c.client.Receieve(c.queueUrl)
			if err != nil {
				//TODO log this entry
				fmt.Println("errro from the queue")
				continue
			}

			if res == nil {
				continue
			}

			err = c.processor.Process(res)
			if err != nil {
				//TODO log this entry
				fmt.Println(err)
			}
			err = c.client.Delete(c.queueUrl, res.ReceiptHandler)
			if err != nil {
				//TODO log this entry
				fmt.Println("errro from the queue")
			}
		}
	}
}
