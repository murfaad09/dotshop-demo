package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type awsConfig struct {
	awsSession *session.Session
}

func InitializeSession(region, accessKeyID, secretAccessKey string) (AWSService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return &awsConfig{awsSession: sess}, nil
}

func (c *awsConfig) GetSession() *session.Session {
	return c.awsSession
}
