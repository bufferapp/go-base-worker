package elasticsearch

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/sha1sum/aws_signing_client"
	elastic "gopkg.in/olivere/elastic.v5"
)

// Client all things ES
type Client struct {
	URL     string
	Client  *elastic.Client
	Context context.Context
}

// NewClient Create a new Elasticsearch client for AWS Elasticsearch Service
// https://github.com/olivere/elastic/wiki/Using-with-AWS-Elasticsearch-Service
func NewClient(awsAccessKeyID string, awsSecretAccessKey string, url string) (*Client, error) {

	ctx := context.Background()

	awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
	signer := v4.NewSigner(awsCredentials)
	awsClient, err := aws_signing_client.New(signer, nil, "es", "us-east-1")
	if err != nil {
		return nil, err
	}

	c, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetScheme("https"),
		elastic.SetHttpClient(awsClient),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		URL:     url,
		Client:  c,
		Context: ctx,
	}, nil
}

// IndexDoc Index document in ES
func (c *Client) IndexDoc(body interface{}, idx string, t string, id string) (bool, error) {
	_, err := c.Client.Index().
		Index(idx).
		Type(t).
		Id(id).
		BodyJson(body).
		Do(c.Context)
	if err != nil {
		// Handle error
		return false, err
	}
	return true, nil
}
