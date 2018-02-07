package queue

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type Queue struct {
	Name    string
	workers map[string]func(interface{})
	conn    *amqp.Connection
}

// NewQueue returns new queue given a name and opens connection server
func NewQueue(name string) (*Queue, error) {
	conn, err := openConnection()
	if err != nil {
		return nil, err
	}
	return &Queue{Name: name, conn: conn, workers: make(map[string]func(interface{}))}, err
}

func (q *Queue) CleanUp() {
	q.conn.Close()
}

// Publish publishes a message to queue
func (q *Queue) Publish(body []byte) error {
	if err := q.validate(); err != nil {
		return err
	}

	ch, err := q.openChannel()
	defer ch.Close()
	if err != nil {
		return err
	}

	return ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
}

// Consume
func (q *Queue) Consume() error {
	if err := q.validate(); err != nil {
		return err
	}
	ch, err := q.openChannel()
	defer ch.Close()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	for message := range msgs {
		go func(m amqp.Delivery) {
			// TODO: make this better
			var args interface{}
			json.Unmarshal(m.Body, &args)
			bodyMap := args.(map[string]interface{})
			name := bodyMap["Name"].(string)
			q.callWorker(name, bodyMap)
			m.Ack(false)
		}(message)
	}

	return nil
}

// RegisterWorker assigns a function to a worker name
func (q *Queue) RegisterWorker(name string, fn func(interface{})) {
	q.workers[name] = fn
}

func (q *Queue) validate() error {
	if q.Name == "" {
		return fmt.Errorf("queue must have name")
	}

	if q.conn == nil {
		return fmt.Errorf("queue must have connection")
	}

	return nil
}

func (q *Queue) callWorker(name string, args interface{}) {
	q.workers[name](args)
}

func openConnection() (*amqp.Connection, error) {
	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = "amqp://localhost"
	}
	return amqp.Dial(url)
}

func (q *Queue) openChannel() (*amqp.Channel, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(
		q.Name, // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func declareQueue(name string, ch *amqp.Channel) (amqp.Queue, error) {
	durable, exclusive := false, false
	autoDelete, noWait := true, true
	return ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, nil)
}
