package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

// Client all things SQS
type Client struct {
	queueURL          string
	client            *sqs.SQS
	visibilityTimeout int64
	maxMessages       int64
	waitTimeSeconds   int64
}

// Option configures the SQS client.
type Option func(*Client)

// WithVisibilityTimeout sets the visibility timeout in seconds for received messages.
// Default is 120 seconds.
func WithVisibilityTimeout(seconds int64) Option {
	return func(c *Client) {
		c.visibilityTimeout = seconds
	}
}

// WithMaxMessages sets the maximum number of messages to receive per poll.
// Valid range is 1-10. Default is 10.
func WithMaxMessages(n int64) Option {
	return func(c *Client) {
		if n < 1 {
			n = 1
		}
		if n > 10 {
			n = 10
		}
		c.maxMessages = n
	}
}

// WithWaitTimeSeconds sets the long-poll wait time in seconds.
// Default is 20 seconds.
func WithWaitTimeSeconds(seconds int64) Option {
	return func(c *Client) {
		c.waitTimeSeconds = seconds
	}
}

// NewClient creates a SQS client. Options can override defaults.
func NewClient(awsAccessKeyID string, awsSecretAccessKey string, queueURL string, opts ...Option) (*Client, error) {

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
		Region:      aws.String("us-east-1"),
	}

	c := &Client{
		queueURL:          queueURL,
		client:            sqs.New(sess, awsConfig),
		visibilityTimeout: 120,
		maxMessages:       10,
		waitTimeSeconds:   20,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Receive receives a single message from the queue.
// For backwards compatibility, returns only the first message even if
// multiple are fetched. Use ReceiveBatch to get all messages.
func (c *Client) Receive() (msg *sqs.Message, err error) {
	msgs, err := c.ReceiveBatch()
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, nil
	}
	return msgs[0], nil
}

// ReceiveBatch receives up to MaxMessages messages from the queue.
func (c *Client) ReceiveBatch() ([]*sqs.Message, error) {
	out, err := c.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: aws.Int64(c.maxMessages),
		VisibilityTimeout:   aws.Int64(c.visibilityTimeout),
		WaitTimeSeconds:     aws.Int64(c.waitTimeSeconds),
	})
	if err != nil {
		return nil, errors.Wrap(err, "receiving sqs message failed")
	}

	return out.Messages, nil
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
