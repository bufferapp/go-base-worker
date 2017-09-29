package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

// Client all things SQS
type Client struct {
	queueURL string
	client   *sqs.SQS
}

// NewClient creates a SQS client.
func NewClient(queueURL string) *Client {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Client{
		queueURL: queueURL,
		client:   sqs.New(sess),
	}
}

// Receive receive a message from the queue.
func (c *Client) Receive() (msg *sqs.Message, err error) {
	var out *sqs.ReceiveMessageOutput
	out, err = c.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(36000),
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		err = errors.Wrap(err, "receiving sqs message failed")
		return
	}

	msg = out.Messages[0]
	return
}

// Delete deletes a message from the queue.
func (c *Client) Delete(msg *sqs.Message) error {
	_, err := c.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		return errors.Wrap(err, "deleting sqs message failed")
	}
	return nil
}
