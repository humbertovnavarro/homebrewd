package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	discordmessage "github.com/humbertovnavarro/homebrewd/pkg/discord-message"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var amqpURI string

const exchangeName = "discord-messages-egress"

type Message struct {
	ChannelID string                 `json:"channel_id"`
	Payload   *discordgo.MessageSend `json:"payload"`
}

func main() {
	godotenv.Load()
	amqpURI = os.Getenv("AMQP_URI")
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_BOT_TOKEN")))
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuilds | discordgo.IntentsMessageContent
	if err != nil {
		logrus.Fatalf("Error creating Discord session: %v", err)
	}
	err = discord.Open()
	if err != nil {
		logrus.Fatalf("Error opening discord session: %v", err)
	}
	defer discord.Close()

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		logrus.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Error creating channel: %v", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logrus.Fatalf("Error declaring exchange: %v", err)
	}

	msgs, err := ch.Consume(
		fmt.Sprintf("%s_queue", exchangeName),
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logrus.Fatalf("Error consuming messages: %v", err)
	}

	go func() {
		for msg := range msgs {
			// Parse message
			var message Message
			if err := json.Unmarshal(msg.Body, &message); err != nil || message.Payload == nil {
				logrus.Printf("Error parsing message: %v", err)
				continue
			}
			fmt.Printf("Received message for channel %s: %v\n", message.ChannelID, message.Payload)
			// Send message to Discord
			_, err := discord.ChannelMessageSendComplex(message.ChannelID, message.Payload)
			if err != nil {
				switch err.Error() {
				case discordmessage.EmptyErrorMessage:
					msg.Reject(false)
				case discordmessage.UnauthorizedErrorMessage:
					msg.Reject(false)
				default:
					msg.Reject(true)
				}
			}
		}
	}()

	logrus.Printf("Listening for messages. Press Ctrl+C to exit.")
	sig := make(chan os.Signal, 1)
	<-sig
}
