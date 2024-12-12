package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Close() error
	Push(msg string) error
}

type queue struct {
	conn    *amqp.Connection
	name    string
	produce *amqp.Channel
}

func New(host, port, name string) (*queue, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port))
	if err != nil {
		return nil, err
	}

	prod, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = prod.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &queue{
		conn:    conn,
		produce: prod,
		name:    name,
	}, nil
}

func (q *queue) Close() error {
	prodErr := q.produce.Close()
	connErr := q.conn.Close()

	if prodErr != nil {
		return prodErr
	}
	if connErr != nil {
		return connErr
	}

	return nil
}

func (q *queue) Push(msg string) error {
	err := q.produce.PublishWithContext(
		context.Background(),
		"",     // exchange
		q.name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		return err
	}

	return nil
}
