package config

import (
	"fmt"
	"github.com/dotunj/pvent/util"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfg Config
)

type Config struct {
	Type     string         `json:"type"`
	Rate     int            `json:"rate"`
	Target   string         `json:"target"`
	Sqs      SqsConfig      `json:"sqs"`
	Google   GoogleConfig   `json:"google"`
	Kafka    KafkaConfig    `json:"kafka"`
	RabbitMq RabbitMqConfig `json:"rabbitmq"`
}

type SqsConfig struct {
	AccessKeyID     string `json:"access_key_id" env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `json:"secret_access_key" env:"AWS_SECRET_ACCESS_KEY"`
	QueueName       string `json:"queue_name" env:"AWS_QUEUE_NAME"`
	Region          string `json:"region" env:"AWS_REGION"`
}

type GoogleConfig struct {
	ProjectID string `json:"project_id" env:"GOOGLE_PROJECT_ID"`
	TopicName string `json:"topic_name" env:"GOOGLE_TOPIC_NAME"`
}

type KafkaConfig struct {
	Address string     `json:"address" env:"KAFKA_ADDRESS"`
	Topic   string     `json:"topic" env:"KAFKA_TOPIC"`
	Auth    *KafkaAuth `json:"auth"`
}

type KafkaAuth struct {
	Type     string `json:"type" env:"KAFKA_AUTH_TYPE"`
	Hash     string `json:"hash" env:"KAFKA_AUTH_HASH"`
	TLS      bool   `json:"tls" env:"KAFKA_AUTH_TLS"`
	Username string `json:"username" env:"KAFKA_AUTH_USERNAME"`
	Password string `json:"password" env:"KAFKA_AUTH_PASSWORD"`
}

type RabbitMqConfig struct {
	Dsn          string `json:"dsn" env:"RABBITMQ_DSN"`
	QueueName    string `json:"queue_name" env:"RABBITMQ_QUEUE_NAME"`
	ExchangeName string `json:"exchange_name" env:"RABBITMQ_EXCHANGE_NAME"`
}

func Load(path string) (*Config, error) {
	_, err := os.Stat(path)

	if err != nil {
		// config file does not exist, read from environment variables
		fmt.Println(">>> READING CONFIG FROM ENV VARIABLES >>>")
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return &cfg, err
		}

		return &cfg, err
	}

	err = cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, err
}

func (c *Config) BindFlags(cmd *cobra.Command) error {
	ptype, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}

	if !util.IsStringEmpty(ptype) {
		cfg.Type = ptype
	}

	rate, err := cmd.Flags().GetInt("rate")
	if err != nil {
		return err
	}

	if rate != 0 {
		cfg.Rate = rate
	}

	target, err := cmd.Flags().GetString("target")
	if err != nil {
		return err
	}

	if !util.IsStringEmpty(target) {
		cfg.Target = target
	}

	return nil
}
