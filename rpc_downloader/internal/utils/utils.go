package utils

import (
	"encoding/json"
	"fmt"
	"rpc/internal/data"
	"github.com/streadway/amqp"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


func ParseJobMessage(body []byte) data.JobMessage {
	//parse the job
	var job data.JobMessage

	err := json.Unmarshal(body, &job)
	if err != nil {
		fmt.Println("error is", err)
	}

	fmt.Printf("Received a job message: %+v", job)

	return job
}

func SendParser(body string)  {
	conn, err := amqp.Dial("amqp://dmx_test:dmx_test@192.168.4.201:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"anker_1_download_results", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	FailOnError(err, "Failed to declare a queue")

	//body := "hello world"
	err = ch.Publish(
		"dmx_test_exchange",     // exchange
		q.Name + "_key", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:  []byte(body),
		})
	FailOnError(err, "Failed to publish a message")
	//log.Printf(" [x] Sent %s", len(body))
}