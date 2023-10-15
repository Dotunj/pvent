package producer

const (
	SQS    = "sqs"
	GOOGLE = "google"
	KAFKA  = "kafka"
)

type Producer interface {
	Broadcast() error
}
