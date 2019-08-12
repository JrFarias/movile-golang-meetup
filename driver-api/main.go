package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
)

//Tracking structure
type Tracking struct {
	PublisherDate time.Time `json:"publisherDate"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	DriverID      string    `json:"driverId"`
}

// SQSService works through the job SQS queue
type SQSService struct {
	Session *session.Session
	SQS     *sqs.SQS
	SQSURL  string
}

func newSQSService(queueURL string, region string) (*SQSService, error) {
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

const queueURL = "http://localhost:4576/queue/DRIVER_TRACKING_QUEUE"
const awsRegion = "us-east-1"

func main() {
	router := gin.Default()
	router.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, time.Now())
	})

	router.POST("/telemetry", func(context *gin.Context) {
		var tracking Tracking
		context.BindJSON(&tracking)
		tracking.PublisherDate = time.Now()
		connect, _ := newSQSService(queueURL, awsRegion)
		result, _ := json.Marshal(tracking)

		sendParams := &sqs.SendMessageInput{
			MessageBody: aws.String(string(result)),
			QueueUrl:    aws.String(queueURL),
		}

		sendResp, err := connect.SQS.SendMessage(sendParams)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("[Send message] \n%v \n\n", sendResp)
	})

	router.Run(":3001")
}
