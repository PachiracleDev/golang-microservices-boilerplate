package rabittmq

import "go.uber.org/fx"

func NewModuleConsumer(service string) fx.Option {
	return fx.Module(
		"rabbitmq_consumer",
		fx.Provide(
			func() ServiceMeta {
				return ServiceMeta{
					ServiceName: service,
				}
			}, NewConsumer),
	)
}

var NewModuleProducer = fx.Module(
	"rabbitmq_producer",
	fx.Provide(NewRabbitMQProducer),
)
