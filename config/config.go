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
