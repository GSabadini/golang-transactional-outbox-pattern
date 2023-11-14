package infra

import (
	"database/sql"
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/producer"
)

type Dependencies struct {
	Broker producer.Broker
	MySQL  *sql.DB
}
