package main

import (
	"github.com/dotunj/pvent/config"
	"github.com/dotunj/pvent/producer"
	"github.com/spf13/cobra"
)

var (
	cfg                    *config.Config
	pType, cfgPath, target string
	rate                   int
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
		p, err := producer.NewProducer(cfg, target)
		if err != nil {
			return err
		}

		return p.Broadcast()
	},
}

func init() {
	dispatcherCmd.Flags().StringVar(&cfgPath, "config", "./pvent.json", "Path to Config file")
	dispatcherCmd.Flags().StringVar(&pType, "type", "", "Message Brokers Type (sqs, google, kafka, rabbitmq)")
	dispatcherCmd.Flags().IntVar(&rate, "rate", 1, "Number of events to send")
	dispatcherCmd.Flags().StringVar(&target, "target", "", "Path to JSON payload to dispatch")
}
