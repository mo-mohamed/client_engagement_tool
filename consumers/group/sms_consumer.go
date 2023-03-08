package groupSmsConsumer

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"customer_engagement/service/queue"
)

type GroupSmsConsumer struct {
	queueUrl    string
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
			// p.ResultsStream <- models.Response{Status: 100, Error: errors.New("cancelled")}
			fmt.Println("Cancelled request")
			return
		default:
			time.Sleep(time.Second * 1)
			sqsClient := newSQS(os.Getenv("AWS_SQS_REGION"), os.Getenv("AWS_SQS_ENDPOINT"))

			request := &sqs.ReceiveMessageInput{
				QueueUrl:              aws.String(os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME")),
				MaxNumberOfMessages:   aws.Int64(1),
				WaitTimeSeconds:       aws.Int64(3),
				MessageAttributeNames: aws.StringSlice([]string{"All"}),
			}
			res, err := queue.GetClient().ReceiveMessage(request)

			fmt.Println(res)

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
			sqsClient.DeleteMessage(d)

		}
	}
}

func newSQS(region, endpoint string) sqsiface.SQSAPI {
	cfg := aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}

	sess := session.Must(session.NewSession(&cfg))
	return sqs.New(sess)
}
