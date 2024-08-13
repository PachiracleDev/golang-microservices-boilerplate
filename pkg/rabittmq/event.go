package rabittmq

import (
	"time"
)

type Event struct {
	Date        time.Time
	Id          string
	AggregateId string
	Payload     interface{}
}
