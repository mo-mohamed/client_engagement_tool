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
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	sqsClient "customer_engagement/clients/sqs"
	"customer_engagement/comm"
	dd "customer_engagement/comm/jobs"
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
			} else {

				// Start processing group - to be refactored
				sqsMessage := res.Messages[0]
				groupId, _ := strconv.Atoi(*sqsMessage.MessageAttributes["GroupId"].StringValue)
				// fmt.Println(groupId)
				// dateString := sqsMessage.MessageAttributes["DateEnqueued"]
				// fmt.Printf("Message: %+v", dateString)
				// fmt.Println()
				dateEn, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", *sqsMessage.MessageAttributes["DateEnqueued"].StringValue)
				// internalId := sqsMessage.MessageAttributes["InternalID"].StringValue
				// fmt.Println(internalId)

				dbHost := os.Getenv("DB_HOST")
				dbUser := os.Getenv("DB_USER")
				dbPassword := os.Getenv("DB_PASSWORD")
				dbPort := os.Getenv("DB_PORT")
				dbName := os.Getenv("DB_NAME")
				dsn := dbUser + ":" + dbPassword + "@tcp" + "(" + dbHost + ":" + dbPort + ")/" + dbName + "?" + "parseTime=true&loc=Local"

				db_conn, _ := sql.Open("mysql", dsn)

				numProfiles := `
					SELECT COUNT(*) FROM group_profile gp JOIN profile p on gp.profile_id = p.id where gp.group_id = ? AND gp.created_at <= ?;
				`
				var count int
				err = db_conn.QueryRow(numProfiles, groupId, dateEn).Scan(&count)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Number of rows are: ", count)
				batchSize := 2
				numOfBatches := count / batchSize

				for i := 0; i < numOfBatches; i++ {
					fmt.Println("--------starting a batch------------")
					var jobs []comm.ICommunication
					for i := 0; i < batchSize; i++ {
						jobs = append(jobs, dd.NewSms(""+strconv.Itoa(i), "hello there"))
					}
					pool := comm.NewComJobPool(2)
					pool.AddBatch(jobs)
					go pool.Run()

					fmt.Println("ok we will wait for messages")
					for i := 0; i < batchSize; i++ {
						<-pool.Results()
					}

					fmt.Println("--------ending a batch------------")
				}
				fmt.Println(numOfBatches)

				// TODO: Process per group

				d := &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME")),
					ReceiptHandle: res.Messages[0].ReceiptHandle,
				}
				awsSqsClient.DeleteMessage(d)
			}
		}
	}
}
