/*
This packahe holds functionality for a group sms consumer.

It starts a number of monitored workers according to the defined number of concurrency.

Keep fetching from the specific SMS groups queue.

Usage:

c := consumer.NewGroupSMSConsumer()
c.Run()

It can be gracefully stopped by calling c.Stop()
*/

package groupSmsConsumer

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

type GroupSmsConsumer struct {
	concurrency int
	done        chan interface{}
	ctx         context.Context
	ctxCancelFn context.CancelFunc
}

func NewGroupSMSConsumer() *GroupSmsConsumer {
	ctx, ctxCancelFn := context.WithCancel(context.TODO())
	concurreny, _ := strconv.Atoi(os.Getenv("GROUP_MS_CONSUMER_CONCURRENCY"))
	return &GroupSmsConsumer{
		concurrency: concurreny,
		done:        make(chan interface{}),
		ctx:         ctx,
		ctxCancelFn: ctxCancelFn,
	}
}

func (consumer *GroupSmsConsumer) Stop() {
	consumer.ctxCancelFn()
}

func (consumer *GroupSmsConsumer) Run() {
	var wg sync.WaitGroup
	for i := 0; i < consumer.concurrency; i++ {
		wg.Add(1)
		go consumer.doWork(&wg, i)

	}

	wg.Wait()
	close(consumer.done)

}

func (consumer *GroupSmsConsumer) doWork(wg *sync.WaitGroup, id int) {
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

			d := &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME")),
				ReceiptHandle: res.Messages[0].ReceiptHandle,
			}
			awsSqsClient.DeleteMessage(d)
		}
	}
}
