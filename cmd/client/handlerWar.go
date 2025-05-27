package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {

	return func(row gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		var ack pubsub.AckType
		outcome, _, _ := gs.HandleWar(row)
		if outcome == gamelogic.WarOutcomeNotInvolved {
			ack = pubsub.NackRequeue
		} else if outcome == gamelogic.WarOutcomeNoUnits {
			ack = pubsub.NackDiscard
		} else if outcome == gamelogic.WarOutcomeOpponentWon {
			ack = pubsub.Ack
		} else if outcome == gamelogic.WarOutcomeYouWon {
			ack = pubsub.Ack
		} else if outcome == gamelogic.WarOutcomeDraw {
			ack = pubsub.Ack
		} else {
			fmt.Println("error: will nackDiscard")
			ack = pubsub.NackDiscard
		}
		return ack
	}
}
