package consumer

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// HandlerFunc is used to define the Handler that is run on for each message
type HandlerFunc func(msg *sqs.Message) error

// HandleMessage is used for the actual execution of each message
func (function HandlerFunc) HandleMessage(msg *sqs.Message) error {
	return function(msg)
}

// Handler interface
type Handler interface {
	HandleMessage(msg *sqs.Message) error
}

// InvalidMessageError for message that can't be processed and should be deleted
type InvalidMessageError struct {
	SQSMessage string
	LogMessage string
}

func (e InvalidMessageError) Error() string {
	return fmt.Sprintf("[Invalid Message: %s] %s", e.SQSMessage, e.LogMessage)
}

// NewInvalidMessageError to create new error for messages that should be deleted
func NewInvalidMessageError(SQSMessage, logMessage string) InvalidMessageError {
	return InvalidMessageError{SQSMessage: SQSMessage, LogMessage: logMessage}
}

// SQSService works through the job SQS queue
type SQSService struct {
	Session *session.Session
	SQS     *sqs.SQS
	SQSURL  string
}

// NewSQSService creates new worker service
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

// Start starts the polling and will continue polling till the application is forcibly stopped
func (service *SQSService) Start(handler Handler) {
	for {
		params := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(service.SQSURL), // Required
			MaxNumberOfMessages: aws.Int64(10),
			MessageAttributeNames: []*string{
				aws.String("All"), // Required
			},
			WaitTimeSeconds: aws.Int64(1),
		}

		resp, err := service.SQS.ReceiveMessage(params)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(resp.Messages) > 0 {
			run(service, handler, resp.Messages)
		}
	}
}

// poll launches goroutine per received message and wait for all message to be processed
func run(service *SQSService, handler Handler, messages []*sqs.Message) {
	numMessages := len(messages)

	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(message *sqs.Message) {
			// launch goroutine
			defer wg.Done()

			if err := handleMessage(service, message, handler); err != nil {
				log.Println(err.Error())
			}
		}(messages[i])
	}

	wg.Wait()
}

func handleMessage(service *SQSService, message *sqs.Message, handler Handler) error {
	err := handler.HandleMessage(message)
	if _, ok := err.(InvalidMessageError); ok {
		log.Println(err.Error())
	} else if err != nil {
		return err
	}

	// Delete the processed (or invalid) message
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(service.SQSURL), // Required
		ReceiptHandle: message.ReceiptHandle,      // Required
	}
	_, err = service.SQS.DeleteMessage(params)
	if err != nil {
		return err
	}
	log.Printf("MessageId %v Processed: \n", message.MessageId)
	return nil
}

// MyError is an error implementation that includes a time and message.
type MyError struct {
	When time.Time
	What string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

//Oops  only the error Message
func Oops() error {
	return MyError{
		time.Now(),
		"the file system has gone away",
	}
}
