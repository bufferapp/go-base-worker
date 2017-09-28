package sqs

import (
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
