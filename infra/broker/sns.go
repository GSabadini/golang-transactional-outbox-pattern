package broker

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNS struct {
	client *sns.Client
}

func NewSNS(config aws.Config) SNS {
	return SNS{client: sns.NewFromConfig(config)}
}

func (s SNS) Publish(ctx context.Context, event any) error {
	jsonMessage, err := json.Marshal(event)
	if err != nil {
		return err
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(jsonMessage)),
		TopicArn: aws.String("arn:aws:sns:sa-east-1:000000000000:Events"),
	}

	_, err = s.client.Publish(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
