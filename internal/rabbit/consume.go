package rabbit

import (
	"context"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func (conn *Conn) StartConsumer(exch, qName, rKey string, handler func(amqp.Delivery) bool) error {
	_, err := conn.Channel.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %v", err)
	}

	err = conn.Channel.QueueBind(qName, "#."+rKey+".#", exch, false, nil)
	if err != nil {
		return fmt.Errorf("queue bind: %v", err)
	}

	// Set prefetchCount above Zero to limit unacknowledged messages
	err = conn.Channel.Qos(0, 0, false)
	if err != nil {
		return err
	}

	// Consumer with explicit ack
	msgs, err := conn.Channel.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %v", err)
	}

	go func() {
		for msg := range msgs {
			if handler(msg) {
				msg.Ack(false)
			} else {
				msg.Nack(false, true)
			}
		}
		log.Fatalf("consumer closed")
	}()

	return nil
}

func (conn *Conn) StartConsumerTemp(ctx context.Context, done chan<- bool, exch, rKey string, handler func(amqp.Delivery) error) error {
	ch, err := conn.Connection.Channel()
	if err != nil {
		return fmt.Errorf("queue declare: %v", err)
	}

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %v", err)
	}

	err = ch.QueueBind(q.Name, "#."+rKey+".#", exch, false, nil)
	if err != nil {
		return fmt.Errorf("queue bind: %V", err)
	}

	// Consume with auto-ack
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %v", err)
	}

	go func() {
		defer ch.Close()
	Consumer:
		for {
			select {
			case msg := <-msgs:
				if err := handler(msg); err != nil {
					done <- true
					break Consumer
				}
			case <-ctx.Done():
				break Consumer
			}
		}
	}()

	return nil
}
