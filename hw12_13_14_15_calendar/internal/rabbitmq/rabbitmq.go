package rabbitmq

import (
	"fmt"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type RabbitClient struct {
	url       string
	conn      *amqp.Connection
	queueName string
	channel   *amqp.Channel
	log       *logger.Logger
}

type Config struct {
	Port int
	Host,
	User,
	Password,
	Exchange,
	Queue string
}

func New(cfg Config, log *logger.Logger) (*RabbitClient, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	rabbitClient := &RabbitClient{
		url: url,
		log: log,
	}

	err := rabbitClient.connect()
	if err != nil {
		return nil, err
	}

	return rabbitClient, nil
}

func (r *RabbitClient) QueueDeclare(queueName string) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.queueName = queueName

	return nil
}

func (r *RabbitClient) Publish(exchangeName string, message []byte) error {
	err := r.channel.Publish(
		exchangeName,
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	return err
}

func (r *RabbitClient) QueueBind(exchangeName string) error {
	return r.channel.QueueBind(
		r.queueName,
		"",
		exchangeName,
		false,
		nil,
	)
}

func (r *RabbitClient) Consume() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitClient) ExchangeDeclare(exchangeName string) error {
	return r.channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}

	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitClient) connect() error {
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	r.conn = conn
	r.channel = channel

	return nil
}
