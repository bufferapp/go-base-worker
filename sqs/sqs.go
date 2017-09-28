package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Type all things SQS
type Type struct {
	queueURL string
	service  *sqs.SQS
}

// Initialize the SQS queue service
func Initialize(url string) Type {

	// TODO: This read ~/.aws/config, Better to explicit the env variable
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	s := Type{
		queueURL: url,
		service:  svc,
	}
	return s
}

// Receive SQS message
func Receive(sqsSvc Type) (*sqs.Message, error) {

	result, err := sqsSvc.service.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sqsSvc.queueURL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(36000),
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		fmt.Println("Something went wrong...", err)
		return nil, err
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return nil, nil
	}

	for _, item := range result.Messages {
		return item, nil
	}

	return nil, nil
}
