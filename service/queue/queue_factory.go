package queue

import (
	"os"
)

func GetClient() *AwsQueueClient {
	if os.Getenv("QUEUE_CLIENT") == "aws" {
		c := newSQS()
		return &c
	}

	return nil
}
