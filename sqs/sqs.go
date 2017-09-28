package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSType all things SQS
type SQSType struct {
	QueueURL string
	Service  *sqs.SQS
}

// Initialize the SQS queue service
func Initialize(url string) SQSType {

	// TODO: This read ~/.aws/config, Better to explicit the env variable
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	s := SQSType{
		QueueURL: url,
		Service:  svc,
	}
	return s
}

// Receive SQS message
func Receive(sqsSvc SQSType) (*sqs.Message, error) {

	result, err := sqsSvc.Service.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sqsSvc.QueueURL),
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

// Delete SQS message with its Handle
func Delete(sqsSvc SQSType, ReceiptHandle *string) bool {

	resultDelete, err := sqsSvc.Service.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(sqsSvc.QueueURL),
		ReceiptHandle: ReceiptHandle,
	})

	if err != nil {
		fmt.Println("Delete Error", err)
		fmt.Println("Delete Output", resultDelete)
		return false
	}

	return true
}
