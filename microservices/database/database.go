package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/DevNiklesh/amqp-redis-messenger-go/internal/config"
	"github.com/DevNiklesh/amqp-redis-messenger-go/internal/models"
	"github.com/DevNiklesh/amqp-redis-messenger-go/internal/rabbit"
	"github.com/streadway/amqp"

	_ "github.com/lib/pq"
)

var conf = config.New()

func main() {
	fmt.Println("[Database service]")

	// Postgres connection
	connPG, err := sql.Open("postgres", conf.PostgresURL+"?sslmode=disable")
	if err != nil {
		log.Fatalf("postgres connection: %s", err)
	}
	defer connPG.Close()

	// Incase table is not created! (This happens when deploying to Cloud Database)
	_, err = connPG.Exec("CREATE TABLE IF NOT EXISTS messages (id SERIAL PRIMARY KEY, message TEXT NOT NULL, created TIMESTAMP NOT NULL)")
	if err != nil {
		log.Fatalf("create table: %s", err)
	}

	// RabbitMQ connection
	connMQ, err := rabbit.GetConn(conf.RabbitURL)
	if err != nil {
		log.Fatalf("rabbit connection: %v", err)
	}
	defer connMQ.Close()

	err = connMQ.DeclareTopicExchange(conf.Exchange)
	if err != nil {
		log.Fatalf("declare exchange: %v", err)
	}

	connMQ.StartConsumer(conf.Exchange, conf.QueueDB, conf.KeyDB, func(d amqp.Delivery) bool {
		return insertToDB(d, connPG)
	})

	select {}
}

func insertToDB(d amqp.Delivery, c *sql.DB) bool {
	var message models.Message
	err := json.Unmarshal(d.Body, &message)
	if err != nil {
		log.Fatalf("unmarshal message: %s", err)
	}

	fmt.Print("message from RabbitMQ: ", message)

	_, err = c.Exec("INSERT INTO messages (message, created) VALUES ($1, to_timestamp($2))", message.Text, message.Time)
	if err != nil {
		log.Fatalf("insert into database: %s", err)
	}

	return true
}
