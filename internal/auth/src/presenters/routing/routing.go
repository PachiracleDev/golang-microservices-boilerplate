package presenters

import (
	"fmt"
	"log"

	"moov/pkg/rabittmq"
	protobuffers "moov/protobuffers"

	"google.golang.org/protobuf/proto"
)

const service = "auth"

func RoutingKey(consumer *rabittmq.Consumer) error {

	msgs := consumer.ConsumeMessages(service)

	forever := make(chan bool)

	go func() {
		log.Printf(" [*] Auth Microservice: Waiting for messages. To exit press CTRL+C")
		for d := range msgs {

			switch d.RoutingKey {
			case "auth.event.signin":

				//DESERIALIZAR REQUEST
				var signInRequest protobuffers.SignInRequest
				proto.Unmarshal(d.Body, &signInRequest)

				//USAR TU CASO DE USO

				//CREAR RESPONSE
				responseBody, _ := proto.Marshal(&protobuffers.SignInResponse{
					AccessToken: "XDD",
				})

				//PUBLISH REPLY
				consumer.PublishReply(
					d,
					responseBody,
				)
			case "auth.event.signup":
				fmt.Println("Message received Signup")

			default:
			}
		}
	}()

	<-forever

	return nil
}
