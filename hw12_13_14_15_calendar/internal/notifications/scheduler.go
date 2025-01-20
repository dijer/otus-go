package notifications

import (
	"encoding/json"
	"time"

	config "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/schedulerconfig"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Notification struct {
	cfg          config.SchedulerConfig
	rabbitClient rabbitmq.RabbitClient
	db           *sqlx.DB
	log          logger.Logger
}

func NewScheduler(
	cfg *config.SchedulerConfig, rabbitClient *rabbitmq.RabbitClient, db *sqlx.DB, log *logger.Logger,
) *Notification {
	return &Notification{
		cfg:          *cfg,
		rabbitClient: *rabbitClient,
		db:           db,
		log:          *log,
	}
}

func (n *Notification) Run() error {
	err := n.rabbitClient.ExchangeDeclare(n.cfg.Rabbit.Exchange)
	if err != nil {
		return err
	}

	err = n.rabbitClient.QueueDeclare(n.cfg.Rabbit.Queue)
	if err != nil {
		return nil
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

	n.log.Info("start scheduler!")
	for {
		select {
		case <-eventsTicker.C:
			notifications, err := n.getNotifications()
			if err != nil {
				n.log.Error(err.Error())
				continue
			}

			err = n.sendNotifications(notifications)
			if err != nil {
				n.log.Error(err.Error())
				continue
			}
		case <-cleanupTicker.C:
			err := n.cleanupOldEvents()
			if err != nil {
				continue
			}
		}
	}
}

func (n *Notification) getNotifications() (*[]storage.Event, error) {
	var events []storage.Event

	query := `
		select id, title, owner, start_time, end_time, description
		from events where start_time > now()
		and notification_sent = false
	`
	err := n.db.Select(&events, query)
	if err != nil {
		return nil, err
	}

	return &events, nil
}

func (n *Notification) sendNotifications(notifications *[]storage.Event) error {
	for _, notification := range *notifications {
		body, err := json.Marshal(notification)
		if err != nil {
			return err
		}

		n.log.Info("scheduler process event", string(notification.ID))

		err = n.rabbitClient.Publish(n.cfg.Rabbit.Exchange, body)
		if err != nil {
			return err
		}

		query := `update events set notification_sent = true where id = $1`
		_, err = n.db.Exec(query, notification.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Notification) cleanupOldEvents() error {
	query := `delete from events where start_time < now() - interval '1 year'`
	_, err := n.db.Exec(query)
	return err
}
