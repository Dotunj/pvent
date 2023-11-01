package producer

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/dotunj/pvent/util"
	"github.com/google/uuid"
	"sync"
	"sync/atomic"
)

type SqsProducer struct {
	rate         int
	svc          *sqs.SQS
	queueURL     *string
	wg           *sync.WaitGroup
	payload      []byte
	errorCount   uint64
	successCount uint64
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
	fmt.Println(">>> Publishing Message via SQS >>>")

	s.wg.Add(s.rate)
	for i := 1; i <= s.rate; i++ {
		go s.dispatch()
	}

	s.wg.Wait()
	return util.PublishReport(s.rate, s.successCount, s.errorCount)
}

func (s *SqsProducer) dispatch() error {
	_, err := s.svc.SendMessage(&sqs.SendMessageInput{
		MessageGroupId:         aws.String(uuid.NewString()),
		MessageDeduplicationId: aws.String(uuid.NewString()),
		QueueUrl:               s.queueURL,
		MessageBody:            aws.String(string(s.payload)),
	})

	defer s.wg.Done()
	defer s.handleError()

	if err != nil {
		fmt.Printf("failed to publish message: %v", err)
		atomic.AddUint64(&s.errorCount, 1)
		return nil
	}

	atomic.AddUint64(&s.successCount, 1)
	return nil
}

func (s *SqsProducer) handleError() {
	if err := recover(); err != nil {
		fmt.Printf("error with sqs client: %v", err)
	}
}
