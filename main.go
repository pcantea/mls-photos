package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/streadway/amqp"
)

const (
	uri       = "amqp://guest:guest@localhost:5672/"
	queueName = "mls-photos"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO - ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR - ", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Println("[*] Connecting to Rabbit...")
	conn, err := amqp.Dial(uri)
	if err != nil {
		// Using fatal so our Docker container can be restarted if there's an issue
		errorLog.Fatalf("Error connecting to RabbitMQ: +%v", err)
	}
	defer conn.Close()
	infoLog.Println("Connected to RabbitMQ successfully.")

	ch, err := conn.Channel()
	if err != nil {
		errorLog.Fatalf("Error opening a channel to RabbitMQ: +%v", err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		errorLog.Fatalf("Error declaring queue %s: +%v", queueName, err)
	}

	msgChan, err := ch.Consume(
		queue.Name,
		"mls-photos",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		errorLog.Fatalf("Error registering consumer with RabbitMQ: +%v", err)
	}

	stopChan := make(chan bool)

	go func() {
		// Do stuff with message here
	}()

	<-stopChan

}

func initAWS(accessKey, secretKey, region string) *session.Session {
	var awsConfig *aws.Config

	if accessKey == "" || secretKey == "" {
		awsConfig = &aws.Config{
			Region: aws.String(region),
		}
	} else {
		awsConfig = &aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		}
	}

	return session.Must(session.NewSession(awsConfig))
}
