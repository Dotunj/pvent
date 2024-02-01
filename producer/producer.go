package producer

import (
	"fmt"
	"strings"

	"github.com/dotunj/pvent/config"
	"github.com/dotunj/pvent/util"
)

const (
	SQS      = "sqs"
	GOOGLE   = "google"
	KAFKA    = "kafka"
	RABBITMQ = "rabbitmq"
)

type Producer interface {
	Broadcast() error
}

func NewProducer(cfg *config.Config, target string) (Producer, error) {
	var p Producer
	var err error

	payload, err := util.ReadJSONTarget(target)
	if err != nil {
		return nil, err
	}

	pType := cfg.Type
	switch {
	case pType == SQS:
		p, err = NewSqsProducer(&SqsConfig{
			AccessKeyID:     cfg.Sqs.AccessKeyID,
			SecretAccessKey: cfg.Sqs.SecretAccessKey,
			Region:          cfg.Sqs.Region,
			QueueName:       cfg.Sqs.QueueName,
			Payload:         payload,
			Rate:            cfg.Rate,
		})
		if err != nil {
			return nil, err
		}
	case pType == GOOGLE:
		p, err = NewGoogleProducer(&GoogleConfig{
			ProjectID: cfg.Google.ProjectID,
			TopicName: cfg.Google.TopicName,
			Payload:   payload,
			Rate:      cfg.Rate,
		})
		if err != nil {
			return nil, err
		}

	case pType == KAFKA:
		kafka := cfg.Kafka
		var auth *KafkaAuth

		if kafka.Auth != nil {
			a := kafka.Auth
			auth = &KafkaAuth{
				Username: a.Username,
				Password: a.Password,
				Type:     a.Type,
				Hash:     a.Hash,
				TLS:      a.TLS,
			}
		}
		p, err = NewKafkaProducer(&KafkaConfig{
			Address: strings.Split(kafka.Address, ","),
			Topic:   kafka.Topic,
			Payload: payload,
			Rate:    cfg.Rate,
			Auth:    auth,
		})

		if err != nil {
			return nil, err
		}

	case pType == RABBITMQ:
		p, err = NewRabbitMqProducer(&RabbitConfig{
			Dsn:       cfg.RabbitMq.Dsn,
			Exchange:  cfg.RabbitMq.ExchangeName,
			QueueName: cfg.RabbitMq.QueueName,
			Payload:   payload,
			Rate:      cfg.Rate,
		})

		if err != nil {
			return nil, err
		}

	default:
		return p, fmt.Errorf("message broker type: %s is not supported", cfg.Type)
	}

	return p, nil
}
