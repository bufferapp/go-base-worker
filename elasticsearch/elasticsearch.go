package elasticsearch

import (
	"context"
	"log"

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
func NewClient(awsAccessKeyID string, awsSecretAccessKey string, url string, env string) (*Client, error) {

	ctx := context.Background()
	var c *elastic.Client
	var err error

	log.Println("Accessing ES cluster ", url)
	if env == "LOCAL" {
		c, err = elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetScheme("http"),
			elastic.SetSniff(false),
		)
	} else {
		awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
		signer := v4.NewSigner(awsCredentials)
		awsClient, err := aws_signing_client.New(signer, nil, "es", "us-east-1")
		if err != nil {
			return nil, err
		}

		c, err = elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetScheme("https"),
			elastic.SetHttpClient(awsClient),
			elastic.SetSniff(false),
		)
	}

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
func (c *Client) IndexDoc(body interface{}, idx string, t string, id string) (*elastic.IndexResponse, error) {
	indexResponse, err := c.Client.Index().
		Index(idx).
		Type(t).
		Id(id).
		BodyJson(body).
		// Refresh("true"). // We want the document available right after the indexing.  refreshing is expensive
		Do(c.Context)
	if err != nil {
		// Handle error
		return indexResponse, err
	}
	return indexResponse, nil
}

// DeleteDoc Delete document in ES
func (c *Client) DeleteDoc(idx string, t string, id string) (*elastic.DeleteResponse, error) {
	deleteResponse, err := c.Client.Delete().
		Index(idx).
		Type(t).
		Id(id).
		Do(c.Context)
	if err != nil {
		// Handle error
		return deleteResponse, err
	}
	return deleteResponse, nil
}
