package notifications

import (
	"encoding/json"
	"fmt"

	notificationcfg "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/notificationconfig"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/streadway/amqp"
)

type Sender struct {
	cfg          notificationcfg.NotificationConfig
	rabbitClient rabbitmq.RabbitClient
	log          logger.Logger
}

func NewSender(
	cfg *notificationcfg.NotificationConfig,
	rabbitClient *rabbitmq.RabbitClient,
	log *logger.Logger,
) *Sender {
	return &Sender{
		cfg:          *cfg,
		rabbitClient: *rabbitClient,
		log:          *log,
	}
}

func (s *Sender) Run() error {
	err := s.rabbitClient.ExchangeDeclare(s.cfg.Rabbit.Exchange)
	if err != nil {
		return err
	}

	err = s.rabbitClient.QueueDeclare(s.cfg.Rabbit.Queue)
	if err != nil {
		return err
	}

	err = s.rabbitClient.QueueBind(s.cfg.Rabbit.Exchange)
	if err != nil {
		return err
	}

	for {
		msgs, err := s.rabbitClient.Consume()
		if err != nil {
			return err
		}

		s.log.Info("sender start!")

		for msg := range msgs {
			err := s.processMessage(msg)
			if err != nil {
				s.log.Error(err.Error())
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		}
	}
}

func (s *Sender) processMessage(msg amqp.Delivery) error {
	var event storage.Event
	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		return err
	}

	s.log.Info(fmt.Sprintf("event sent: %v", event))

	return nil
}
