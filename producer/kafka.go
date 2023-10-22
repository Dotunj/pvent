package producer

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"strconv"
	"sync"
	"time"
)

type KafkaProducer struct {
	rate      int
	address   []string
	topic     string
	partition int
	writer    *kafka.Writer
	wg        *sync.WaitGroup
	payload   []byte
	auth      *KafkaAuth
	ctx       context.Context
}

type KafkaConfig struct {
	Address []string
	Topic   string
	Rate    int
	Payload []byte
	Auth    *KafkaAuth
}

type KafkaAuth struct {
	Type     string
	Hash     string
	TLS      bool
	Username string
	Password string
}

func NewKafkaProducer(cfg *KafkaConfig) (*KafkaProducer, error) {
	sharedTransport, err := getTransport(cfg.Auth)
	if err != nil {
		return nil, err
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Address...),
		Topic:                  cfg.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		Transport:              sharedTransport,
	}

	return &KafkaProducer{
		address: cfg.Address,
		topic:   cfg.Topic,
		rate:    cfg.Rate,
		writer:  writer,
		payload: cfg.Payload,
		wg:      &sync.WaitGroup{},
		ctx:     context.Background(),
	}, nil

}

func getTransport(auth *KafkaAuth) (*kafka.Transport, error) {
	var mechanism sasl.Mechanism
	var err error

	if auth == nil {
		return nil, nil
	}

	if auth.Type != "plain" && auth.Type != "scram" {
		return nil, fmt.Errorf("auth type: %s is not supported", auth.Type)
	}

	if auth.Type == "plain" {
		mechanism = plain.Mechanism{
			Username: auth.Username,
			Password: auth.Password,
		}
	}

	if auth.Type == "scram" {
		algo := scram.SHA512

		if auth.Hash == "SHA256" {
			algo = scram.SHA256
		}

		mechanism, err = scram.Mechanism(algo, auth.Username, auth.Password)
		if err != nil {
			return nil, err
		}
	}

	sharedTransport := &kafka.Transport{
		SASL:        mechanism,
		DialTimeout: 15 * time.Second,
	}

	if auth.TLS {
		sharedTransport.TLS = &tls.Config{}
	}

	return sharedTransport, nil
}

func (k *KafkaProducer) Broadcast() error {
	k.wg.Add(k.rate)

	for i := 1; i <= k.rate; i++ {
		go k.dispatch(i)
	}

	k.wg.Wait()

	return nil
}

func (k *KafkaProducer) dispatch(i int) error {
	ctx, cancel := context.WithTimeout(k.ctx, 10*time.Second)
	defer cancel()

	defer k.wg.Done()

	err := k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.Itoa(i)),
		Value: k.payload,
	})

	if err != nil {
		return fmt.Errorf("unable to publish message - %v", err)
	}

	return nil
}
