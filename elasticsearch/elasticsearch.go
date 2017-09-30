package elasticsearch

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/sha1sum/aws_signing_client"
	elastic "gopkg.in/olivere/elastic.v5"
)

// NewClient Create a new Elasticsearch client for AWS Elasticsearch Service
// https://github.com/olivere/elastic/wiki/Using-with-AWS-Elasticsearch-Service
func NewClient(awsAccessKeyID string, awsSecretAccessKey string, url string) (*elastic.Client, error) {

	awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
	signer := v4.NewSigner(awsCredentials)
	awsClient, err := aws_signing_client.New(signer, nil, "es", "us-east-1")
	if err != nil {
		return nil, err
	}

	return elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetScheme("https"),
		elastic.SetHttpClient(awsClient),
		elastic.SetSniff(false),
	)
}
