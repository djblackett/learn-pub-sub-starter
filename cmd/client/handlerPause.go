package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {

	return func(state routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		if gs.Paused {
			gs.HandlePause(routing.PlayingState{
				IsPaused: false,
			})
		} else {
			gs.HandlePause(routing.PlayingState{
				IsPaused: true,
			})
		}
		return pubsub.Ack
	}
}
