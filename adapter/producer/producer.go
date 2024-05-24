package producer

import (
	"context"
	"encoding/json"
	"github.com/GSabadini/golang-transactional-outbox-pattern/domain"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type (
	Producer struct {
		client *sns.Client
		topic  string
	}
)

func NewProducer(client *sns.Client, topic string) Producer {
	return Producer{
		client: client,
		topic:  topic,
	}
}

func (p Producer) Publish(ctx context.Context, event domain.Event) error {
	ctx, span := opentelemetry.NewSpan(ctx, "producer.producer.publish")
	defer span.End()

	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(p.topic),
		Message:  aws.String(string(message)),
	}

	_, err = p.client.Publish(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (p Producer) PublishBatch(ctx context.Context, events []domain.Event) error {
	ctx, span := opentelemetry.NewSpan(ctx, "producer.producer.publish_batch")
	defer span.End()

	var batchEntries []types.PublishBatchRequestEntry

	for _, event := range events {
		message, err := json.Marshal(event)
		if err != nil {
			return err
		}

		batchEntries = append(batchEntries, types.PublishBatchRequestEntry{
			Id:      aws.String(uuid.NewString()),
			Message: aws.String(string(message)),
		})
	}

	input := &sns.PublishBatchInput{
		PublishBatchRequestEntries: batchEntries,
		TopicArn:                   aws.String(p.topic),
	}

	_, err := p.client.PublishBatch(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
