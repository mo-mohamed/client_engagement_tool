/*
This packahe holds functionality for a group sms consumer.

It starts a number of monitored workers according to the defined number of concurrency.

Keep fetching from the specific SMS groups queue.

Usage:

c := consumer.NewGroupSMSConsumer()
c.Run()

It can be gracefully stopped by calling c.Stop()
*/

package groupQueueConsumer

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	sqsClient "customer_engagement/clients/sqs"
)

type GroupQueueConsumer struct {
	concurrency int
	done        chan interface{}
	ctx         context.Context
	ctxCancelFn context.CancelFunc
}

func NewGroupQueueConsumer() *GroupQueueConsumer {
	ctx, ctxCancelFn := context.WithCancel(context.TODO())
	concurreny, _ := strconv.Atoi(os.Getenv("GROUP_MS_CONSUMER_CONCURRENCY"))
	return &GroupQueueConsumer{
		concurrency: concurreny,
		done:        make(chan interface{}),
		ctx:         ctx,
		ctxCancelFn: ctxCancelFn,
	}
}

func (consumer *GroupQueueConsumer) Stop() {
	consumer.ctxCancelFn()
}

func (consumer *GroupQueueConsumer) Run() {
	var wg sync.WaitGroup
	for i := 0; i < consumer.concurrency; i++ {
		wg.Add(1)
		go consumer.doWork(&wg, i)

	}

	wg.Wait()
	close(consumer.done)

}

func (consumer *GroupQueueConsumer) doWork(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for {
		select {

		case <-consumer.ctx.Done():
			return
		default:
			time.Sleep(time.Second * 1)
			awsSqsClient := sqsClient.New(os.Getenv("AWS_SQS_REGION"), os.Getenv("AWS_SQS_ENDPOINT"))

			request := &sqs.ReceiveMessageInput{
				QueueUrl:              aws.String(os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME")),
				MaxNumberOfMessages:   aws.Int64(1),
				WaitTimeSeconds:       aws.Int64(3),
				MessageAttributeNames: aws.StringSlice([]string{"All"}),
			}

			res, err := awsSqsClient.ReceiveMessage(request)

			if err != nil {
				fmt.Println("Error from sms queue: ", err)
				continue
			}

			if len(res.Messages) == 0 {
				continue
			}

			// TODO: Process per group
			// countProfiles := 10
			// batshsize := 2
			// batches := countProfiles / batshsize
			// skip := 0
			// for i := 0; i < batches; i++ {
			// 	fmt.Println("Starting a btahc")
			// 	var jobs []comm.ICommunication
			// 	for i := 0; i < batshsize; i++ {
			// 		jobs = append(jobs, dd.NewSms(""+strconv.Itoa(i), "hello there"))
			// 	}
			// 	pool := comm.NewComJobPool(10)
			// 	pool.AddBatch(jobs)
			// 	pool.Run()
			// 	skip = skip + batshsize
			// 	for elem := range pool.Results() {
			// 		fmt.Println(elem)
			// 	}
			// 	fmt.Println("Batch ended")
			// }

			// fmt.Println("--------Finalized Batch----------")

			d := &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME")),
				ReceiptHandle: res.Messages[0].ReceiptHandle,
			}
			awsSqsClient.DeleteMessage(d)
		}
	}
}
