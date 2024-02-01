package producer

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/dotunj/pvent/util"
	"sync"
	"sync/atomic"
	"time"
)

type GoogleProducer struct {
	rate         int
	client       *pubsub.Client
	wg           *sync.WaitGroup
	topic        *pubsub.Topic
	payload      []byte
	ctx          context.Context
	errorCount   uint64
	successCount uint64
}

type GoogleConfig struct {
	ProjectID string
	TopicName string
	Rate      int
	Payload   []byte
}

func NewGoogleProducer(cfg *GoogleConfig) (Producer, error) {
	client, err := pubsub.NewClient(context.Background(), cfg.ProjectID)
	if err != nil {
		fmt.Println("error with setting up client", err)
		return nil, err
	}

	g := &GoogleProducer{
		rate:    cfg.Rate,
		client:  client,
		wg:      &sync.WaitGroup{},
		topic:   client.Topic(cfg.TopicName),
		payload: cfg.Payload,
		ctx:     context.Background(),
	}

	return g, nil
}

func (g *GoogleProducer) Broadcast() error {
	fmt.Println(">>> Publishing Message via Google Pub/Sub >>>")

	g.wg.Add(g.rate)
	for i := 1; i <= g.rate; i++ {
		go g.dispatch(i)
	}

	g.wg.Wait()
	return util.PublishReport(g.rate, g.successCount, g.errorCount)
}

func (g *GoogleProducer) dispatch(i int) error {
	ctx, cancel := context.WithTimeout(g.ctx, 10*time.Second)
	defer cancel()
	defer g.wg.Done()
	defer g.handleError()

	result := g.topic.Publish(ctx, &pubsub.Message{
		Data: g.payload,
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(i int, res *pubsub.PublishResult) {
		defer wg.Done()

		_, err := res.Get(ctx)
		if err != nil {
			fmt.Printf("failed to publish message: %v", err)
			atomic.AddUint64(&g.errorCount, 1)
			return
		}

		atomic.AddUint64(&g.successCount, 1)
	}(i, result)
	wg.Wait()

	return nil
}

func (g *GoogleProducer) handleError() {
	if err := recover(); err != nil {
		fmt.Printf("error with google client: %v", err)
	}
}
