package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/movile-golang-meetup/driver-store-service/src/config"
	"github.com/movile-golang-meetup/driver-store-service/src/consumer"
)

//Tracking structure
type Tracking struct {
	PublisherDate time.Time `json:"publisherDate"`
	DriverID      string    `json:"driverId"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	TrackingTime  time.Time `json:"trackingTime"`
	ID            string    `json:"id"`
}

// DB var used to connect with PG
var DB *sql.DB

func saveOnDatabase(tracking Tracking) Tracking {
	tracking.TrackingTime = time.Now()
	var trackid string

	err := DB.QueryRow(`INSERT INTO driver_trackings(driver_id, created, publisher_time, latitude, longitude)
	VALUES($1,$2, $3, $4, $5) RETURNING id`, tracking.DriverID, tracking.TrackingTime, tracking.PublisherDate, tracking.Latitude, tracking.Longitude).Scan(&trackid)
	tracking.ID = trackid

	if len(trackid) == 0 {
		log.Fatalf("error %v", err)
	}

	return tracking
}

func sendToOrderSQS(tracking Tracking) {
	const queueURL = "http://localhost:4576/queue/DRIVER_QUEUE"
	const awsRegion = "us-east-1"

	connect, _ := config.NewSQSService(queueURL, awsRegion)
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
}

func processSQS(msg *sqs.Message) error {
	if msg.Body != nil {
		var tracking Tracking
		json.Unmarshal([]byte(*msg.Body), &tracking)

		result := saveOnDatabase(tracking)
		sendToOrderSQS(tracking)
		log.Printf("result %v", result)

		return nil
	}

	return consumer.Oops()
}

func main() {
	worker, err := consumer.NewSQSService("http://localhost:4576/queue/DRIVER_TRACKING_QUEUE", "us-east-1")

	if err != nil {
		log.Printf("Error creating new Worker Service: %s\n", err)
	}

	dbconnection := config.DBConnect()
	DB = dbconnection

	log.Println("consumer running...")
	worker.Start(consumer.HandlerFunc(processSQS))
}
