package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/DevNiklesh/amqp-redis-messenger-go/internal/config"
	"github.com/DevNiklesh/amqp-redis-messenger-go/internal/models"
	"github.com/DevNiklesh/amqp-redis-messenger-go/internal/rabbit"
	"github.com/streadway/amqp"
)

var conf = config.New()

func main() {
	fmt.Println("[Backend service]")

	// RabbitMQ connection
	conn, err := rabbit.GetConn(conf.RabbitURL)
	if err != nil {
		log.Fatalf("rabbit connection: %s", err)
	}
	defer conn.Close()

	// Declaring the exchange topic
	err = conn.DeclareTopicExchange(conf.Exchange)
	if err != nil {
		log.Fatalf("declare exchange: %s", err)
	}

	conn.StartConsumer(conf.Exchange, conf.QueueBack, conf.KeyBack, printMessages)

	publishInput(conn)
}

func printMessages(d amqp.Delivery) bool {
	var message models.Message
	err := json.Unmarshal(d.Body, &message)
	if err != nil {
		log.Fatalf("unmarshal message: %s", err)
	}

	fmt.Printf("> %s\n", string(message.Text))

	return true
}

func publishInput(c *rabbit.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		inputTime := time.Now().Unix() // in milliseconds
		inputMsg := models.Message{Text: input, Source: "back", Time: inputTime}
		message, err := json.Marshal(inputMsg)
		if err != nil {
			log.Fatalf("marshal message: %s", err)
		}

		key := conf.KeyBack + "." + conf.KeyDB
		err = c.Publish(conf.Exchange, key, message)
		if err != nil {
			log.Fatalf("publish message: %s", err)
		}
	}
}
