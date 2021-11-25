package rabbit

import (
	"crypto/tls"

	"github.com/streadway/amqp"
)

type Conn struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func GetConn(rabbitUrl string) (*Conn, error) {
	cfg := new(tls.Config)
	cfg.InsecureSkipVerify = true

	conn, err := amqp.DialTLS(rabbitUrl, cfg)
	if err != nil {
		return &Conn{}, err
	}

	ch, err := conn.Channel()
	return &Conn{
		Connection: conn,
		Channel:    ch,
	}, err
}

func (conn *Conn) Close() error {
	return conn.Connection.Close()
}
