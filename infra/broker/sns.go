package broker

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewSNS(config aws.Config) *sns.Client {
	return sns.NewFromConfig(config)
}
