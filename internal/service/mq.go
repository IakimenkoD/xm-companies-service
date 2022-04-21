package service

import (
	"encoding/json"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

//go:generate minimock -i MessageQueue -g -o mq_mock.go
type MessageQueue interface {
	NotifyCompanyUpdated(company *model.Company) error
}

type messageQueue struct {
	channel *amqp.Channel
	queue   amqp.Queue
	log     *zap.Logger
}

type NotificationTask struct {
	CompanyID int64      `json:"company_id"`
	UpdatedAt *time.Time `json:"updated_at"`

	NewName    *string `json:"new_name,omitempty"`
	NewCode    *string `json:"new_code,omitempty"`
	NewCountry *string `json:"new_country,omitempty"`
	NewWebsite *string `json:"new_website,omitempty"`
	NewPhone   *string `json:"new_phone,omitempty"`
}

func (m messageQueue) NotifyCompanyUpdated(company *model.Company) error {
	task := NotificationTask{
		CompanyID:  company.ID,
		NewName:    &company.Name,
		NewCode:    &company.Code,
		NewCountry: &company.Country,
		NewWebsite: &company.Website,
		NewPhone:   &company.Phone,
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return errors.Wrap(err, "marshalling notification task")
	}
	m.log.Debug("publishing notification event to mq", zap.ByteString("event", taskBytes))

	return m.channel.Publish(
		"",
		m.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        taskBytes,
		})
}

func NewMessageQueue(cfg *config.Config, log *zap.Logger) (MessageQueue, error) {
	conn, err := amqp.Dial(cfg.MQ.Address)
	if err != nil {
		return nil, err
	}
	amqpChannel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := amqpChannel.QueueDeclare(
		cfg.MQ.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &messageQueue{
		channel: amqpChannel,
		queue:   queue,
		log:     log,
	}, nil
}
