package infra

import (
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Dependencies struct {
	SNS   *sns.Client
	MySQL *sql.DB
}
