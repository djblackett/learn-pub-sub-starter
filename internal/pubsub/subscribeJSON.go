package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {

	ch, _, err := DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, simpleQueueType)
	if err != nil {
		return err
	}

	newChan, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for message := range newChan {
			var decoded T
			err = json.Unmarshal(message.Body, &decoded)
			ack := handler(decoded)

			if ack == Ack {
				err = message.Ack(false)
			} else if ack == NackDiscard {
				err = message.Nack(false, false)
			} else if ack == NackRequeue {
				err = message.Nack(false, true)
			}
			fmt.Println("ack-type:", ack)

			if err != nil {
				fmt.Println("error:", err)
			}
		}
	}()

	return nil
}
