package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Println("Connection successful")

	input, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	username := input

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, "pause."+input, routing.PauseKey, 1)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	gameState := gamelogic.NewGameState(input)

	// handler := handlerPause(gameState)
	pauseHandler := handlerPause(gameState)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, "pause."+input, routing.PauseKey, 1, pauseHandler)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	// army_moves
	ch, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+input, routing.ArmyMovesPrefix+".*", 1)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	moveHandler := handlerMove(gameState, ch)
	err = pubsub.SubscribeJSON(conn, string(routing.ExchangePerilTopic), routing.ArmyMovesPrefix+"."+input, routing.ArmyMovesPrefix+".*", 1, moveHandler)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	for {
		input := gamelogic.GetInput()

		if len(input) == 0 {
			continue
		} else if input[0] == "spawn" {
			err = gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		} else if input[0] == "move" {
			mv, err := gameState.CommandMove(input)
			fmt.Println(mv)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = pubsub.PublishJSON(ch, string(routing.ExchangePerilTopic), string(routing.ArmyMovesPrefix)+"."+username, mv)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Move successful")
			fmt.Println("Move published")
		} else if input[0] == "status" {
			gameState.CommandStatus()
		} else if input[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if input[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		} else if input[0] == "quit" {
			gamelogic.PrintQuit()
		} else {
			fmt.Println("Unknown command")
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
