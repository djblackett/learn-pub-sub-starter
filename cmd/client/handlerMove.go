package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	defer fmt.Print(">")
	var ack pubsub.AckType
	return func(am gamelogic.ArmyMove) pubsub.AckType {
		outcome := gs.HandleMove(am)
		gs.GetPlayerSnap()
		if outcome == gamelogic.MoveOutComeSafe {
			ack = pubsub.Ack
		} else if outcome == gamelogic.MoveOutcomeMakeWar {
			err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix+"."+gs.GetUsername(), outcome)
			if err != nil {
				fmt.Println("error publishing event/message")
			}
			ack = pubsub.NackRequeue
		} else {
			ack = pubsub.NackDiscard
		}
		return ack
	}
}
