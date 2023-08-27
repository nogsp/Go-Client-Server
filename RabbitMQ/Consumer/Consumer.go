package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// configurations
const RequestQueue = "request_queue"
const ResponseQueue = "response_queue"

type Request struct {
	Num int
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func fibo(n int) int {
	ans := 1
	prev := 0
	for i := 1; i < n; i++ {
		temp := ans
		ans = ans + prev
		prev = temp
	}
	return ans
}

func main() {
	//connect to broker
	conn, err := amqp.Dial("amqp://guest:guest@172.17.0.2:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//create channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//create request queue
	q, err := ch.QueueDeclare(
		RequestQueue, //routing key(queue's name)
		false,        //durable
		false,        //autodelete
		false,        //exclusive
		false,        //nowait
		nil,          //args
	)
	failOnError(err, "Failed to create request queue(broker's queue)")

	//create request queue's consumer
	msgs, err := ch.Consume(
		q.Name, //routing key(queue's name)
		"",     //consumer
		true,   //autoACK
		false,  //exclusive
		false,  //noLocal
		false,  //nowait
		nil,    //args
	)
	failOnError(err, "Failed to create response queue's consumer(broker's consumer)")

	fmt.Println("Consumer is ready!")

	//receive and process the messages
	for d := range msgs {

		//receive request
		msg := Request{}
		err := json.Unmarshal(d.Body, &msg)
		failOnError(err, "Failed to deserialize the message")

		//process request
		replymsgBytes, err := json.Marshal(fibo(msg.Num))
		fmt.Println(fibo(msg.Num), msg.Num)
		failOnError(err, "Failed to serialize the message")

		//send(publish) response
		err = ch.Publish(
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          replymsgBytes,
			})
		failOnError(err, "Failed to send the message to the broker")
	}
}
