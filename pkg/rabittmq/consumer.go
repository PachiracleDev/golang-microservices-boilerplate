package rabittmq

import (
	"fmt"
	"moov/config"
	"moov/pkg/errors"

	"github.com/rabbitmq/amqp091-go"
)

type ServiceMeta struct {
	ServiceName string
}

type Consumer struct {
	Channel *amqp091.Channel
	Conn    *amqp091.Connection
}

func NewConsumer(conf *config.Config, metaData ServiceMeta) *Consumer {
	conn, _ := GetConnection(conf.RabbitMQ.Uri)
	channel, _ := conn.Channel()

	err := channel.QueueBind(
		fmt.Sprintf("%s-queue", metaData.ServiceName),   // queue name
		fmt.Sprintf("%s.event.*", metaData.ServiceName), // routing key
		conf.RabbitMQ.Exchange,                          // exchange
		false,
		nil,
	)
	if err != nil {
		errors.WrapErrorf(err, errors.ErrorCodeUnknown, "channel.QueueBind")
		return nil
	}

	return &Consumer{
		Channel: channel,
		Conn:    conn,
	}
}

func (u *Consumer) PublishReply(d amqp091.Delivery, responseBody []byte) {
	u.Channel.Publish(
		"",        // Exchange donde se publicará la respuesta
		d.ReplyTo, // Routing key de la respuesta
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType:   "application/x-protobuf",
			Body:          responseBody,
			CorrelationId: d.CorrelationId,
		},
	)
}

func (u *Consumer) ConsumeMessages(service string) <-chan amqp091.Delivery {
	msgs, err := u.Channel.Consume(
		fmt.Sprintf("%s-queue", service),    // queue
		fmt.Sprintf("%s-consumer", service), // consumer
		true,                                // auto-ack
		false,                               // exclusive
		false,                               // no-local
		false,                               // no-wait
		nil,                                 // args
	)

	if err != nil {
		return nil
	}

	return msgs
}
