package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

		serverUrl := os.Getenv("AMQP_SERVER")
		routingKey := os.Getenv("AMQP_ROUTING_KEY")
		if routingKey == "" {
			routingKey = "webhooks"
		}
		conn, err := amqp.Dial(serverUrl)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		b, err := ioutil.ReadAll(r.Body)

		type RequestFormat struct {
			Method string
			Url    *url.URL
			Headers map[string][]string
			Body   string
		}

		group := RequestFormat{
			Method: r.Method,
			Url:    r.URL,
			Headers: r.Header,
			Body:   string(b),
		}

		body, err := json.Marshal(&group)

		err = ch.Publish(
			"",         // exchange
			routingKey, // routing key
			false,      // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})

		failOnError(err, "Failed to publish a message")

}

func main() {

	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	fmt.Println("Number of CPUs: ", nCPU)
	http.HandleFunc("/", handler)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}

}
