package discordmessage

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/streadway/amqp"
)

type Sender struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewSender(amqpURI, exchangeName, exchangeType string) (*Sender, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Sender{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (s *Sender) Send(exchange, routingKey string, message *discordgo.MessageSend) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling discord message: %v", err)
	}
	if messageJSON == nil {
		return fmt.Errorf("error marshalling discord message: %s", "message is nil")
	}
	return s.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageJSON,
		},
	)
}

func (s *Sender) Close() error {
	return s.conn.Close()
}
