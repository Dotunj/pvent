package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	cfg Config
)

type Config struct {
	Type   string       `json:"type"`
	Rate   int          `json:"rate"`
	Target string       `json:"target"`
	Sqs    SqsConfig    `json:"sqs"`
	Google GoogleConfig `json:"google"`
	Kafka  KafkaConfig  `json:"kafka"`
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

func Load(path string) (Config, error) {
	_, err := os.Stat(path)

	// config file does not exist
	// read from env instead
	if err != nil {
		fmt.Println(">>> Reading from Environment variables >>>>")
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return cfg, err
		}
	} else {
		// reading from config file
		fmt.Println(">>> Reading from Config file >>>>")
		err := cleanenv.ReadConfig(path, &cfg)
		if err != nil {
			return cfg, err
		}
	}

	return cfg, err
}

func (c *Config) BindFlags(cmd *cobra.Command) error {
	ptype, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}

	if !IsStringEmpty(ptype) {
		cfg.Type = ptype
	}

	accessKeyID, err := cmd.Flags().GetString("access-key-id")
	if err != nil {
		return err
	}

	if !IsStringEmpty(accessKeyID) {
		cfg.Sqs.AccessKeyID = accessKeyID
	}

	secretAccessKey, err := cmd.Flags().GetString("secret-access-key")
	if err != nil {
		return err
	}

	if !IsStringEmpty(secretAccessKey) {
		cfg.Sqs.SecretAccessKey = secretAccessKey
	}

	queueName, err := cmd.Flags().GetString("queue-name")
	if err != nil {
		return err
	}

	if !IsStringEmpty(queueName) {
		cfg.Sqs.QueueName = queueName
	}

	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return err
	}

	if !IsStringEmpty(region) {
		cfg.Sqs.Region = region
	}

	projectID, err := cmd.Flags().GetString("project-id")
	if err != nil {
		return err
	}

	if !IsStringEmpty(projectID) {
		cfg.Google.ProjectID = projectID
	}

	gtopic, err := cmd.Flags().GetString("gtopic")
	if err != nil {
		return err
	}

	if !IsStringEmpty(gtopic) {
		cfg.Google.TopicName = gtopic
	}

	ktopic, err := cmd.Flags().GetString("ktopic")
	if err != nil {
		return err
	}

	if !IsStringEmpty(ktopic) {
		cfg.Kafka.Topic = ktopic
	}

	address, err := cmd.Flags().GetString("address")
	if err != nil {
		return err
	}

	if !IsStringEmpty(address) {
		cfg.Kafka.Address = address
	}

	authType, err := cmd.Flags().GetString("auth")
	if err != nil {
		return err
	}

	if !IsStringEmpty(authType) {
		cfg.Kafka.Auth.Type = authType
	}

	hash, err := cmd.Flags().GetString("hash")
	if err != nil {
		return err
	}

	if !IsStringEmpty(hash) {
		cfg.Kafka.Auth.Hash = hash
	}

	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return err
	}

	if !IsStringEmpty(username) {
		cfg.Kafka.Auth.Username = username
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return err
	}

	if !IsStringEmpty(password) {
		cfg.Kafka.Auth.Password = password
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

	if !IsStringEmpty(target) {
		cfg.Target = target
	}

	return nil
}

func IsStringEmpty(s string) bool { return len(strings.TrimSpace(s)) == 0 }
