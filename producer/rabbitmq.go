package producer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dotunj/pvent/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Exchange string

type RabbitProducer struct {
	channel  *amqp.Channel
	queue    amqp.Queue
	exchange string
	// Payload message to send
	payload []byte
	// Number of messages to send
	rate         int
	wg           *sync.WaitGroup
	errorCount   uint64
	successCount uint64
}

const (
	DIRECT  Exchange = "direct"
	TOPIC   Exchange = "topic"
	HEADERS Exchange = "headers"
	FANOUT  Exchange = "fanout"
)

type RabbitConfig struct {
	Rate      int
	Payload   []byte
	Exchange  string
	Dsn       string
	QueueName string
}

func NewRabbitMqProducer(cfg *RabbitConfig) (Producer, error) {
	conn, err := amqp.Dial(cfg.Dsn)
	if err != nil {
		fmt.Printf("failed to connect to RabbitMQ: %s", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("failed to open channel: %s", err)
		return nil, err
	}

	q, err := ch.QueueDeclare(cfg.QueueName, false, false, false, false, nil)
	if err != nil {
		fmt.Printf("failed to declare queue: %s", err)
		return nil, err
	}

	r := &RabbitProducer{
		queue:   q,
		channel: ch,
		payload: cfg.Payload,
		rate:    cfg.Rate,
		wg:      &sync.WaitGroup{},
	}

	return r, nil
}

func (r *RabbitProducer) Broadcast() error {
	fmt.Println(">>> Publishing Message via RabbitMQ >>>")

	r.wg.Add(r.rate)
	for i := 1; i <= r.rate; i++ {
		go r.dispatch()
	}

	r.wg.Wait()
	return util.PublishReport(r.rate, r.successCount, r.errorCount)
}

func (r *RabbitProducer) dispatch() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(
		ctx,
		r.exchange,
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        r.payload,
		},
	)

	if err != nil {
		fmt.Printf("failed to publish message: %v", err)
		atomic.AddUint64(&r.errorCount, 1)
		return nil
	}

	atomic.AddUint64(&r.successCount, 1)
	return nil
}
