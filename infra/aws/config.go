package aws

import (
	"context"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewConfig(ctx context.Context) (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   env.AWSProvider,
			URL:           env.AWSEndpoint,
			SigningRegion: env.AWSRegion,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(env.AWSRegion),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}
