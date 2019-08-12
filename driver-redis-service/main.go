package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/movile-golang-meetup/driver-redis-service/config"
	"github.com/movile-golang-meetup/driver-redis-service/consumer"
)

//Tracking structure
type Tracking struct {
	OriginalEventDate string    `json:"originalEventDate"`
	OrderID           string    `json:"orderId"`
	Eta               uint64    `json:"eta"`
	PublisherDate     time.Time `json:"publisherDate"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	TrackingTime      time.Time `json:"trackingTime"`
	ID                string    `json:"id"`
}

func saveOnRedis(tracking Tracking) {
	const addr = "localhost:6379"
	const pass = ""

	result, _ := json.Marshal(tracking)
	err := config.RedisConnect(addr, pass).Set("DRIVE_TRACKING", result, 0).Err()
	if err != nil {
		log.Println("Success to save on redis!")
	}

}

func process(msg *sqs.Message) error {
	if msg.Body != nil {
		var tracking Tracking
		json.Unmarshal([]byte(*msg.Body), &tracking)
		saveOnRedis(tracking)

		return nil
	}

	return consumer.Oops()
}

func main() {
	worker, err := consumer.NewSQSService("http://localhost:4576/queue/DRIVER_QUEUE", "us-east-1")
	log.Println("consumer running...")

	if err != nil {
		log.Printf("Error creating new Worker Service: %s\n", err)
	}

	worker.Start(consumer.HandlerFunc(process))

}
