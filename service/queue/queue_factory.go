package queue

import (
	"os"
)

func GetClient() IGroupQueue {
	if os.Getenv("QUEUE_CLIENT") == "aws" {
		return &AwsQueueClient{}
	}

	return nil
}
