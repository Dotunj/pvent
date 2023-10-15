package producer

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"sync"
)

type SqsProducer struct {
	rate     int
	svc      *sqs.SQS
	queueURL *string
	wg       *sync.WaitGroup
	payload  []byte
}

type SqsConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	QueueName       string
	Payload         []byte
	Rate            int
}

func NewSqsProducer(cfg *SqsConfig) (Producer, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})

	if err != nil {
		fmt.Println("error with initializing session", err)
		return nil, err
	}

	//create a service client
	svc := sqs.New(sess)
	url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(cfg.QueueName),
	})

	if err != nil {
		fmt.Println("failed to fetch queue url", err)
		return nil, err
	}

	s := &SqsProducer{
		rate:     cfg.Rate,
		svc:      svc,
		queueURL: url.QueueUrl,
		wg:       &sync.WaitGroup{},
		payload:  cfg.Payload,
	}

	return s, nil
}

func (s *SqsProducer) Broadcast() error {
	s.wg.Add(s.rate)

	for i := 1; i <= s.rate; i++ {
		go s.dispatch()
	}

	s.wg.Wait()

	return nil
}

func (s *SqsProducer) dispatch() error {
	_, err := s.svc.SendMessage(&sqs.SendMessageInput{
		MessageGroupId:         aws.String(uuid.NewString()),
		MessageDeduplicationId: aws.String(uuid.NewString()),
		QueueUrl:               s.queueURL,
		MessageBody:            aws.String(string(s.payload)),
	})

	defer s.wg.Done()

	if err != nil {
		return err
	}

	fmt.Println("Message sent successfully")
	return nil
}
