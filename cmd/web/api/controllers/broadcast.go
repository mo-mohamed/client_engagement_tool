package controllers

import (
	vmodels "customer_engagement/cmd/web/api/view_models"
	dbc "customer_engagement/data_store/config"
	dbm "customer_engagement/data_store/models"
	repository "customer_engagement/data_store/repository"
	groupProducer "customer_engagement/producers/message"
	queue "customer_engagement/queue"
	sqsClient "customer_engagement/queue/awssqs"
	"fmt"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type BroadcastController struct{}

func (BroadcastController) BroadcastGroup() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var bcr vmodels.BroadcastRequest
		json.NewDecoder(r.Body).Decode(&bcr)
		ok, errors := bcr.Validate()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}

		groupId := bcr.GroupId
		gRepo := repository.NewRepository[dbm.Group](dbc.DB)

		dbGroup, exists := gRepo.Exists(groupId)
		if !exists {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		_, err := produceGroup(dbGroup.ID, bcr.MessageBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			fmt.Println("SUCCEDDED")
		}

	}
}

func produceGroup(groupId int, message string) (string, error) {
	cfg := aws.Config{
		Region:   aws.String(os.Getenv("AWS_SQS_REGION")),
		Endpoint: aws.String(os.Getenv("AWS_SQS_ENDPOINT")),
	}

	sess := session.Must(session.NewSession(&cfg))
	client := sqsClient.NewSqs(sqs.New(sess))
	s := groupProducer.NewGroupProducer(client)

	attributes := make([]queue.Attribute, 0)
	attributes = append(attributes, queue.Attribute{
		Key:   "GroupId",
		Value: strconv.Itoa(groupId),
		Type:  "String",
	})
	attributes = append(attributes, queue.Attribute{
		Key:   "DateEnqueued",
		Value: time.Now().UTC().String(),
		Type:  "String",
	})
	attributes = append(attributes, queue.Attribute{
		Key:   "InternalID",
		Value: uuid.NewString(),
		Type:  "String",
	})

	messageRequest := queue.SendRequest{
		QueueUrl:   os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Body:       message,
		Attributes: attributes,
	}
	return s.Produce(&messageRequest)
}
