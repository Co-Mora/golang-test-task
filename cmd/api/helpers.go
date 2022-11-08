package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

type requestHeader struct {
	ContentType string `header:"ContentType"`
}

func (app *Config) readJSON(c *gin.Context) error {

	connectRabbitMQ, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	msgs, err := channelRabbitMQ.Consume(
		"MessageService",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Recieved Message: %s\n", d.Body)
		}
	}()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever

	return nil
}

func (app *Config) writeJSON(c *gin.Context, status int, data any, headers ...http.Header) {

	connectRabbitMQ, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	out, err := json.Marshal(data)

	h := requestHeader{
		ContentType: "application/json",
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	message := amqp.Publishing{
		ContentType:  "application/json",
		Body:         out,
		DeliveryMode: amqp.Persistent,
	}

	// Attempt to publish a message to the queue.
	if err := channelRabbitMQ.Publish(
		"",
		"MessageService",
		false,
		false,
		message,
	); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, "Successfully Sent")
}
