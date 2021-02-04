package main

import (
	"log"
	"math/rand"
	"os"
	"rpc/internal/data"
	"strings"
	"time"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func callRPC(job string) (res string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(job),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {

			res := string(d.Body)
			return res
			}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// get shop id and request to rpc server
	Shop_id := bodyFrom(os.Args)

	job := fmt.Sprintf("{\"shop_id\": %v, \"interval\": 600}", Shop_id)

	log.Printf(" [x] Requesting job: (%s)", job)

	res := callRPC(job)

	var shopPage data.ShopPage

	err := json.Unmarshal([]byte(res), &shopPage)
	if err != nil {
		log.Println("error is", err)
	}

	log.Printf(" [.] Got %+v", shopPage)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] ==	 "" {
		s = "16461019"
	} else {
		s = strings.Join(args[1:], " ")
	}

	return s
}