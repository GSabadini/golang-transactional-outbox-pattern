package env

import "time"

const (
	ServerPort            = ":8080"
	ServerShutdownTimeout = 30 * time.Second

	AWSProvider   = "aws"
	AWSRegion     = "sa-east-1"
	AWSEndpoint   = "http://localhost:4566"
	SNSEventTopic = "arn:aws:sns:sa-east-1:000000000000:Events"

	TracerEndpoint   = "localhost:4317"
	TracerGlobalName = "golang-transactional-outbox-pattern"

	CronInterval = 30 * time.Second

	DBDriver   = "mysql"
	DBUser     = "dev"
	DBPassword = "dev"
	DBEndpoint = "localhost"
	DBName     = "dev"
)
