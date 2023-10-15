package main

import (
	"github.com/dotunj/pvent/config"
	"github.com/dotunj/pvent/producer"
	"github.com/dotunj/pvent/util"
	"github.com/spf13/cobra"
)

var (
	cfg                                                                     config.Config
	pType, cfgPath, accessKeyID, secretAccessKey, queueName, region, target string
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
			})

			if err != nil {
				return err
			}
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
	dispatcherCmd.Flags().IntVar(&rate, "rate", 1, "Number of events to send")
	dispatcherCmd.Flags().StringVar(&target, "target", "", "Path to JSON payload to dispatch")

	rootCmd.AddCommand(dispatcherCmd)
}
