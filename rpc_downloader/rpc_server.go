package main

import (
	"encoding/json"
	"fmt"
	"log"
	"rpc/internal/data"
	"rpc/internal/utils"
	"bytes"
	"github.com/streadway/amqp"
	"net/http"
	"sync"
	"time"
)

func crawlUrl(url string) string {
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")

	res, getErr := client.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()

	return newStr
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			// parse the Job
			job := utils.ParseJobMessage(d.Body)

			log.Println("\nshop id: ", job.Shop_id)

			url := fmt.Sprintf("https://shopee.vn/api/v2/search_items/?by=pop&limit=100&match_id=%v&newest=0&order=desc&page_type=shop&version=2", job.Shop_id)

			res := crawlUrl(url)

			var shopPage data.ShopPage

			err := json.Unmarshal([]byte(res), &shopPage)
			if err != nil {
				log.Println("error is", err)
			}

			// limit the number of goroutine
			maxGoroutines := 10
			guard := make(chan int, maxGoroutines)

			var wg sync.WaitGroup

			// Crawl detail product
			numProduct := len(shopPage.Items)
			log.Println("number of product", numProduct)

			for i := 0; i < numProduct; i++ {
				wg.Add(1)
				guard <- i
				go func(i int) {
					defer wg.Done()
					itemId := shopPage.Items[i].Itemid

					productUrl := fmt.Sprintf("https://shopee.vn/api/v2/item/get?itemid=%v&shopid=16461019", itemId)

					// crawl the detail of item and send to parser
					itemInfo := crawlUrl(productUrl)
					utils.SendParser(itemInfo)

					var product data.DetailProduct

					err := json.Unmarshal([]byte(itemInfo), &product)
					if err != nil {
						fmt.Println("error is", err)
					}

					log.Printf("crawl product with id: %+v", product.Item.Itemid)

					<-guard
				}(i)
			}

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(res),
				})
			utils.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}