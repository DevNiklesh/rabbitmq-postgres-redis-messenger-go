package rabbit

import (
	"time"

	"github.com/streadway/amqp"
)

func (conn Conn) Publish(exch, rKey string, message []byte) error {
	return PublishInChannel(conn.Channel, exch, rKey, message)
}

func PublishInChannel(ch *amqp.Channel, exch, rKey string, message []byte) error {
	return ch.Publish(
		exch,  // exchange name
		rKey,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         message,
		},
	)
}
