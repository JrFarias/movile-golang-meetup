package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSService works through the job SQS queue
type SQSService struct {
	Session *session.Session
	SQS     *sqs.SQS
	SQSURL  string
}

// NewSQSService used to connect with sqs service
func NewSQSService(queueURL string, region string) (*SQSService, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		MaxRetries:  aws.Int(5),
		Credentials: credentials.NewStaticCredentials("foo", "var", ""),
		Endpoint:    aws.String(queueURL),
	})
	svc := sqs.New(sess)

	builder := &SQSService{
		Session: sess,
		SQS:     svc,
		SQSURL:  queueURL,
	}

	return builder, nil
}
