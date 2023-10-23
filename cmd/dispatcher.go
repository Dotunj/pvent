package main

import (
	"fmt"
	"github.com/dotunj/pvent/config"
	"github.com/dotunj/pvent/producer"
	"github.com/dotunj/pvent/util"
	"github.com/spf13/cobra"
	"strings"
)

var (
	cfg                                                                     *config.Config
	pType, cfgPath, accessKeyID, secretAccessKey, queueName, region, target string
	projectID, gtopic, ktopic, address, authType, hash, username, password  string
	rate                                                                    int
)

var dispatcherCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "dispatches events across message brokers",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load(cfgPath)
		if err != nil {
			return err
		}

		err = cfg.BindFlags(cmd)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var dispatcher producer.Producer
		var err error

		payload, err := util.ReadJSONTarget(target)
		if err != nil {
			return err
		}

		if cfg.Type == producer.SQS {
			sqs := cfg.Sqs
			dispatcher, err = producer.NewSqsProducer(&producer.SqsConfig{
				AccessKeyID:     sqs.AccessKeyID,
				SecretAccessKey: sqs.SecretAccessKey,
				Region:          sqs.Region,
				QueueName:       sqs.QueueName,
				Payload:         payload,
				Rate:            cfg.Rate,
			})

			if err != nil {
				return err
			}
		}

		if cfg.Type == producer.GOOGLE {
			google := cfg.Google
			dispatcher, err = producer.NewGoogleProducer(&producer.GoogleConfig{
				ProjectID: google.ProjectID,
				TopicName: google.TopicName,
				Payload:   payload,
				Rate:      cfg.Rate,
			})

			if err != nil {
				return err
			}
		}

		if cfg.Type == producer.KAFKA {
			kafka := cfg.Kafka
			var auth *producer.KafkaAuth

			if kafka.Auth != nil {
				a := kafka.Auth
				auth = &producer.KafkaAuth{
					Username: a.Username,
					Password: a.Password,
					Type:     a.Type,
					Hash:     a.Hash,
					TLS:      a.TLS,
				}
			}
			dispatcher, err = producer.NewKafkaProducer(&producer.KafkaConfig{
				Address: strings.Split(kafka.Address, ","),
				Topic:   kafka.Topic,
				Payload: payload,
				Rate:    cfg.Rate,
				Auth:    auth,
			})
		}

		if dispatcher == nil {
			return fmt.Errorf("message broker type: %s is not supported", cfg.Type)
		}

		return dispatcher.Broadcast()
	},
}

func init() {
	dispatcherCmd.Flags().StringVar(&cfgPath, "config", "./pvent.json", "Path to Config file")
	dispatcherCmd.Flags().StringVar(&pType, "type", "", "Message Brokers Type (sqs, google, kafka)")
	dispatcherCmd.Flags().StringVar(&accessKeyID, "access-key-id", "", "AWS Access Key ID")
	dispatcherCmd.Flags().StringVar(&secretAccessKey, "secret-access-key", "", "AWS Secret Access Key")
	dispatcherCmd.Flags().StringVar(&queueName, "queue-name", "", "AWS Queue Name")
	dispatcherCmd.Flags().StringVar(&region, "region", "", "AWS Region")
	dispatcherCmd.Flags().StringVar(&projectID, "project-id", "", "Pub/Sub Project ID")
	dispatcherCmd.Flags().StringVar(&gtopic, "gtopic", "", "Pub/Sub Topic Name")
	dispatcherCmd.Flags().StringVar(&ktopic, "ktopic", "", "Kafka Topic Name")
	dispatcherCmd.Flags().StringVar(&address, "address", "", "Kafka Cluster Address")
	dispatcherCmd.Flags().StringVar(&authType, "auth", "", "Kafka Auth Type")
	dispatcherCmd.Flags().StringVar(&hash, "hash", "", "Kafka Auth Hash")
	dispatcherCmd.Flags().StringVar(&username, "username", "", "Kafka Auth Username")
	dispatcherCmd.Flags().StringVar(&password, "password", "", "Kafka Auth Password")
	dispatcherCmd.Flags().IntVar(&rate, "rate", 1, "Number of events to send")
	dispatcherCmd.Flags().StringVar(&target, "target", "", "Path to JSON payload to dispatch")

	rootCmd.AddCommand(dispatcherCmd)
}
