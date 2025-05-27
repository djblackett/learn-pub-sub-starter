package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Println("Connection successful")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	defer ch.Close()
	fmt.Println("Channel created successfully")

	pauseData := routing.PlayingState{
		IsPaused: false,
	}

	resumeData := routing.PlayingState{
		IsPaused: true,
	}

	pauseBytes, err := json.Marshal(pauseData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resumeBytes, err := json.Marshal(resumeData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, string(routing.PauseKey), pauseBytes)
	// if err != nil {
	// 	fmt.Println("error", err)
	// 	os.Exit(1)
	// }

	// err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.GameLogSlug+".*", "hi")
	// if err != nil {
	// 	fmt.Println("error", err)
	// 	os.Exit(1)
	// }

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, "game_logs", routing.GameLogSlug, 0)

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		} else if input[0] == "pause" {
			fmt.Println("Sending pause message")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, string(routing.PauseKey), pauseBytes)
			if err != nil {
				fmt.Println("error", err)
				os.Exit(1)
			}
		} else if input[0] == "resume" {
			fmt.Println("Resuming game")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, string(routing.PauseKey), resumeBytes)
			if err != nil {
				fmt.Println("error", err)
				os.Exit(1)
			}
		} else if input[0] == "quit" {
			fmt.Println("Exiting program...")
			break
		} else {
			fmt.Println("Sorry, unknown command")
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down")
}
