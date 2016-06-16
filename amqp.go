package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type AmqpConnection struct {
	conn *amqp.Connection
}

func createAmqpConnection(user string, password string, host string) (*AmqpConnection, error) {

	uri := "amqp://" + user + ":" + password + "@" + host
	c := &AmqpConnection{conn: nil}

	var err error

	log.Printf("Dialing to %q", uri)
	c.conn, err = amqp.Dial(uri)

	if err != nil {
		return nil, fmt.Errorf("Dial error to %q: %s", uri, err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	return c, nil

}

func (c *AmqpConnection) Shutdown() error {

	defer log.Printf("AMQP Connection shutdown OK")

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	return nil

}

func PublishMesssage(host string, user string, password string, exchange string, key string, mandatory bool, waitConfirm bool, msg *amqp.Publishing) error {

	amqpConnection, err := createAmqpConnection(user, password, host)

	if err != nil {
		return err
	}

	channel, err := amqpConnection.conn.Channel()

	if err != nil {
		return fmt.Errorf("Channel error: %s", err)
	}

	var confirms <-chan amqp.Confirmation

	if waitConfirm {
		confirms = channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	err = channel.Publish(exchange, key, mandatory, false, *msg)

	if err != nil {
		return fmt.Errorf("Publish error: %s", err)
	}

	var confirmed amqp.Confirmation

	if waitConfirm {
		if confirmed = <-confirms; confirmed.Ack {
			fmt.Printf("Published, exchange: %s, key: %s, message: %q", exchange, key, msg)
		} else {
			//log.Printf("Not published to AMQP, %q, %d", delivery.CorrelationId, confirmed.DeliveryTag)
			return fmt.Errorf("Not published, exchange: %s, key: %s, message: %q", exchange, key, msg)
		}
	}

	err = amqpConnection.Shutdown()

	if err != nil {
		return err
	}

	return nil
}
