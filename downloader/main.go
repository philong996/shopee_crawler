package main

import (
	"fmt"
	//"io/ioutil"
	"sync"
	"log"
	"net/http"
	//"os"
	"encoding/json"
	"bytes"
	"time"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessage(body string)  {
	conn, err := amqp.Dial("amqp://dmx_test:dmx_test@192.168.4.201:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"anker_1_download_results", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")

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
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", len(body))
}

func Crawl(url string, wg *sync.WaitGroup) {
	//defer wg.Done()

	shopeeClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")

	res, getErr := shopeeClient.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()

	sendMessage(newStr)
}

type JobMessage struct {
	Url string
	Interval int
}

func parseJobMessage(body []byte) JobMessage {
	//parse the job
	var job JobMessage

	err := json.Unmarshal(body, &job)
	if err != nil {
		fmt.Println("error is", err)
	}

	log.Printf("Received a job message: %s", job)

	return job
}

func main()  {
	conn, err := amqp.Dial("amqp://dmx_test:dmx_test@192.168.4.201:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"anker_1", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
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

	var wg sync.WaitGroup

	forever := make(chan bool)

	maxDownloadRoutine := 10
	for i := 0; i < maxDownloadRoutine; i++ {
		wg.Add(1)
		go func(i int) {
			for {
				fmt.Println(i, "\n")
				defer wg.Done()

				d := <- msgs

				//parse job message
				job := parseJobMessage(d.Body)

				// crawl the url
				Crawl(job.Url, &wg)
			}
		}(i)
	}


	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}