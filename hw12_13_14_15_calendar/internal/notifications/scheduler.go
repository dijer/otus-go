package notifications

import (
	"context"
	"encoding/json"
	"time"

	notificationcfg "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/notificationconfig"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Notification struct {
	cfg          notificationcfg.NotificationConfig
	rabbitClient rabbitmq.RabbitClient
	storage      storage.Storage
	log          logger.Logger
}

func NewScheduler(
	cfg *notificationcfg.NotificationConfig,
	rabbitClient *rabbitmq.RabbitClient,
	storage storage.Storage,
	log *logger.Logger,
) *Notification {
	return &Notification{
		cfg:          *cfg,
		rabbitClient: *rabbitClient,
		storage:      storage,
		log:          *log,
	}
}

func (n *Notification) Run(ctx context.Context) error {
	err := n.rabbitClient.ExchangeDeclare(n.cfg.Rabbit.Exchange)
	if err != nil {
		return err
	}

	err = n.rabbitClient.QueueDeclare(n.cfg.Rabbit.Queue)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(n.cfg.Scheduler.Interval)
	if err != nil {
		return err
	}

	eventsTicker := time.NewTicker(duration)
	defer eventsTicker.Stop()

	cleanupDuration, err := time.ParseDuration(n.cfg.Scheduler.Cleanup)
	if err != nil {
		return err
	}
	cleanupTicker := time.NewTicker(cleanupDuration)
	defer cleanupTicker.Stop()

	n.log.Info("start scheduler!")
	for {
		select {
		case <-eventsTicker.C:
			notifications, err := n.storage.GetNotifications(ctx)
			if err != nil {
				n.log.Error(err.Error())
				continue
			}

			err = n.sendNotifications(ctx, notifications)
			if err != nil {
				n.log.Error(err.Error())
				continue
			}
		case <-cleanupTicker.C:
			err := n.storage.CleanupOldEvents(ctx)
			if err != nil {
				n.log.Error(err.Error())
				continue
			}
		}
	}
}

func (n *Notification) sendNotifications(ctx context.Context, notifications []storage.Event) error {
	for _, notification := range notifications {
		body, err := json.Marshal(notification)
		if err != nil {
			return err
		}

		n.log.Info("scheduler process event", string(notification.ID))

		err = n.rabbitClient.Publish(n.cfg.Rabbit.Exchange, body)
		if err != nil {
			return err
		}

		err = n.storage.SendNotification(ctx, notification.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
