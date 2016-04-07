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

	if r.Method == "POST" {
		serverUrl := os.Getenv("AMQP_SERVER")
		conn, err := amqp.Dial(serverUrl)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		b, err := ioutil.ReadAll(r.Body)

		type RequestFormat struct {
			method string
			url    *url.URL
			headers map[string][]string
			body   string
		}
		group := RequestFormat{
			method: r.Method,
			url:    r.URL,
			headers: r.Header,
			body:   string(b),
		}
		body, err := json.Marshal(&group)

		err = ch.Publish(
			"",         // exchange
			"webhooks", // routing key
			false,      // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError(err, "Failed to publish a message")
	}

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
