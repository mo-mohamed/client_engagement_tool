package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

func main() {

	go func() {
		x := "sss"
		fmt.Println(x)
	}()

	fmt.Println("ook")
	fmt.Println("ook2")

	// sqsClient := newSQS(os.Getenv("AWS_SQS_REGION"), os.Getenv("AWS_SQS_ENDPOINT"))
	// _, err := sendMessage(sqsClient, "12345", string(os.Getenv("AWS_SQS_ENDPOINT"))+"/"+os.Getenv("AWS_SQS_SMS_GROUP_NAME"))
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"))
	// fmt.Println(os.Getenv("AWS_SQS_ENDPOINT"))

	// attrs := make([]queue.Attribute, 0)
	// attrs = append(attrs, queue.Attribute{
	// 	Key:   "test",
	// 	Value: "val",
	// 	Type:  "String",
	// })

	// req := queue.SendRequest{
	// 	QueueURL:   os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
	// 	Body:       "HE~llo this is a body",
	// 	Attributes: attrs,
	// }

	// s := producers.GroupMessageProducer{}
	// s.EnqueueGroupBroadcast()(req)
	// time.Sleep(time.Second * 2)

	// con := groupSmsConsumer.NewGroupSMSConsumer()
	// // fun := func() {
	// // 	time.Sleep(time.Second * 3)
	// // 	con.Stop()
	// // }
	// // go fun()
	// con.Run()

	// fmt.Println(os.Getenv("AWS_SQS_REGION"))
	// sqsClient := newSQS(os.Getenv("AWS_SQS_REGION"), os.Getenv("AWS_SQS_ENDPOINT"))
	// res, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
	// 	QueueUrl:              aws.String(os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME")),
	// 	MaxNumberOfMessages:   aws.Int64(10),
	// 	WaitTimeSeconds:       aws.Int64(20),
	// 	MessageAttributeNames: aws.StringSlice([]string{"All"}),
	// })

	// fmt.Println(res)
	// fmt.Println(err)

}

func newSQS(region, endpoint string) sqsiface.SQSAPI {
	cfg := aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}

	sess := session.Must(session.NewSession(&cfg))
	return sqs.New(sess)
}
