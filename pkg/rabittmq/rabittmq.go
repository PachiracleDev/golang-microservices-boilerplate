package rabittmq

import (
	"fmt"
	"log"
	"moov/config"
	"moov/pkg/errors"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ ...
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Exchange   string
}

func NewRabbitMQProducer(conf *config.Config) *RabbitMQ {

	conn, _ := GetConnection(conf.RabbitMQ.Uri)

	channel, err := conn.Channel()
	if err != nil {
		return nil
	}

	err = channel.ExchangeDeclare(
		conf.RabbitMQ.Exchange, // name
		"topic",                // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return nil
	}

	//CREAR LAS COLAS
	for _, service := range conf.RabbitMQ.Services {
		_, err := channel.QueueDeclare(
			fmt.Sprintf("%s-queue", service), // name
			true,                             // durable
			false,                            // delete when unused
			false,                            // exclusive
			false,                            // no-wait
			nil,                              // arguments
		)
		if err != nil {
			return nil
		}

		err = channel.QueueBind(
			fmt.Sprintf("%s-queue", service), // queue name
			fmt.Sprintf("%s.event.*", service),
			conf.RabbitMQ.Exchange, // exchange
			false,
			nil,
		)
		if err != nil {
			return nil
		}
	}

	//ESTO SIRVE PARA QUE EL CONSUMIDOR NO RECIBA MAS DE UN MENSAJE A LA VEZ Y SE LO PASE A OTRO CONSUMIDOR
	if err := channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil
	}

	// Dead Letter Exchange Proximamente :v

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
		Exchange:   conf.RabbitMQ.Exchange,
	}
}

func GetConnection(uri string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, errors.WrapErrorf(err, errors.ErrorCodeUnknown, "amqp.Dial")
	}

	log.Println("Connected to RabbitMQ")

	return conn, nil
}

// Close ...
func (r *RabbitMQ) Close() {
	r.Connection.Close()
}

type SendEvent struct {
	Event      []byte
	Service    string
	RoutingKey string
}

func (r *RabbitMQ) Send(body SendEvent) error {

	err := r.Channel.Publish(
		r.Exchange,
		fmt.Sprintf("%s.event.%s", body.Service, body.RoutingKey),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/x-protobuf",
			Body:        body.Event,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return errors.WrapErrorf(err, errors.ErrorCodeUnknown, "ch.Publish")
	}

	return nil
}

func (r *RabbitMQ) SendAndListen(body SendEvent) ([]byte, error) {
	// Configura la cola de respuesta temporal
	replyQueue, err := r.Channel.QueueDeclare(
		"",    // Nombre de la cola (vacío para una cola temporal)
		false, // durable
		true,  // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to encode event: %v", err)
	}

	correlationId := uuid.New().String()

	// Publica el mensaje con la cabecera ReplyTo
	err = r.Channel.Publish(
		r.Exchange, // Exchange
		fmt.Sprintf("%s.event.%s", body.Service, body.RoutingKey), // Routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "application/x-protobuf",
			Body:          body.Event,
			Timestamp:     time.Now(),
			CorrelationId: correlationId,   // Correlation ID único
			ReplyTo:       replyQueue.Name, // Cola donde esperas la respuesta
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return nil, err
	}

	// Consume la respuesta desde la cola de respuesta
	msgs, err := r.Channel.Consume(
		replyQueue.Name, // Nombre de la cola
		"",              // Consumer
		true,            // Auto Ack
		false,           // Exclusive
		false,           // No Local
		false,           // No Wait
		nil,             // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Espera la respuesta y devuelve al cliente HTTP
	for d := range msgs {
		if d.CorrelationId == correlationId {
			return d.Body, nil
		}
	}

	return nil, nil
}
