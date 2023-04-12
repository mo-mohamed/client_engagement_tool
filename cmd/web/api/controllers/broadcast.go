package controllers

import (
	viewModels "customer_engagement/cmd/web/api/view_models"
	dbconfig "customer_engagement/data_store/config"
	db_models "customer_engagement/data_store/models"
	repository "customer_engagement/data_store/repository"
	Queue "customer_engagement/queue"
	sqsClient "customer_engagement/queue/awssqs"
	"fmt"
	"os"
	"strconv"
	"time"

	"encoding/json"

	newprod "customer_engagement/producers/message"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
)

type BroadcastController struct{}

func (BroadcastController) BroadcastGroup() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var bcr viewModels.BroadcastRequest
		json.NewDecoder(r.Body).Decode(&bcr)
		ok, errors := bcr.Validate()
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
			return
		}

		groupId := bcr.GroupId
		gRepo := repository.NewRepository[db_models.Group](dbconfig.DB)

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
	client := sqsClient.NewSqs(sess)
	s := newprod.NewGroupProducer(client)

	attributes := make([]Queue.Attribute, 0)
	attributes = append(attributes, Queue.Attribute{
		Key:   "GroupId",
		Value: strconv.Itoa(groupId),
		Type:  "String",
	})
	attributes = append(attributes, Queue.Attribute{
		Key:   "DateEnqueued",
		Value: time.Now().UTC().String(),
		Type:  "String",
	})
	attributes = append(attributes, Queue.Attribute{
		Key:   "InternalID",
		Value: uuid.NewString(),
		Type:  "String",
	})

	messageRequest := Queue.SendRequest{
		QueueUrl:   os.Getenv("AWS_SQS_ENDPOINT") + "/" + os.Getenv("AWS_SQS_SMS_GROUP_NAME"),
		Body:       message,
		Attributes: attributes,
	}
	return s.Produce(&messageRequest)
}
