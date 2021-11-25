package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DevNiklesh/amqp-redis-messenger-go/internal/config"
	"github.com/DevNiklesh/amqp-redis-messenger-go/internal/rabbit"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

var (
	conf = config.New()
	ctx  = context.Background()
)

func main() {
	fmt.Println("[Cache service]")

	// Redis connection
	connR := redis.NewClient(&redis.Options{
		Addr:     conf.RedisURL,
		Password: "",
		DB:       0,
	})

	// RabbitMQ connection
	connMQ, err := rabbit.GetConn(conf.RabbitURL)
	if err != nil {
		log.Fatalf("rabbit connection: %s", err)
	}
	defer connMQ.Close()

	err = connMQ.DeclareTopicExchange(conf.Exchange)
	if err != nil {
		log.Fatalf("declare exchange: %s", err)
	}

	connMQ.StartConsumer(conf.Exchange, conf.QueueCache, conf.KeyCache, func(d amqp.Delivery) bool {
		return updateRedis(d, connR)
	})

	select {}
}

func updateRedis(d amqp.Delivery, c *redis.Client) bool {
	fmt.Println("Received msg from RabbitMQ")

	// Add a message, limit to 10 in cache, increment total count
	if _, err := c.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.LPush(ctx, "messages", d.Body)
		p.LTrim(ctx, "messages", 0, 9)
		p.Incr(ctx, "total")
		return nil
	}); err != nil {
		log.Fatalf("update redis: %s", err)
	}

	return true
}
