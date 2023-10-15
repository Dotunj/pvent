package producer

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type GoogleProducer struct {
	rate        int
	client      *pubsub.Client
	wg          *sync.WaitGroup
	topic       *pubsub.Topic
	payload     []byte
	ctx         context.Context
	totalErrors uint64
}

type GoogleConfig struct {
	ProjectID string
	TopicName string
	Target    string
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
	for i := 1; i <= g.rate; i++ {
		g.dispatch(i)
	}

	if g.totalErrors > 0 {
		return fmt.Errorf("%d of %d messages did not publish successfully", g.totalErrors, g.rate)
	}

	return nil
}

func (g *GoogleProducer) dispatch(i int) error {
	ctx, cancel := context.WithTimeout(g.ctx, 10*time.Second)
	defer cancel()

	result := g.topic.Publish(ctx, &pubsub.Message{
		Data: g.payload,
	})

	g.wg.Add(1)
	go func(i int, res *pubsub.PublishResult) {
		defer g.wg.Done()

		_, err := res.Get(ctx)
		if err != nil {
			fmt.Println("failed to publish", err)
			atomic.AddUint64(&g.totalErrors, 1)
		}

	}(i, result)

	g.wg.Wait()

	return nil
}
