package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Other configurations
const SampleSize = 10000
const RequestQueue = "request_queue"
const ResponseQueue = "response_queue"

type Message struct {
	Num int
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}

func calculate_mean(l *list.List) float64 {
	total := 0.0
	for e := l.Front(); e != nil; e = e.Next() {
		if num, ok := e.Value.(float64); ok {
			total += num
		}
	}
	return float64(total / float64(l.Len()))
}

func main() {
	//connect to broker
	conn, err := amqp.Dial("amqp://guest:guest@172.17.0.3:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//create channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//create response queue
	replyQueue, err := ch.QueueDeclare(
		ResponseQueue+"/"+strconv.Itoa(os.Getpid()), //routing key(queue's name)
		false, //durable
		false, //autodelete
		true,  //exclusive
		false, //nowait
		nil,   //args
	)
	failOnError(err, "Failed to create response queue")

	//create response queue's consumer
	msgs, err := ch.Consume(
		replyQueue.Name, //routing key(queue's name)
		"",              //consumer
		true,            //autoACK
		false,           //exclusive
		false,           //noLocal
		false,           //nowait
		nil,             //args
	)
	failOnError(err, "Failed to create response queue's consumer")

	arr_times := list.New()

	fmt.Println("Producer is ready!")

	//send the message
	for i := 0; i < SampleSize; i++ {
		msg := Message{
			Num: i,
		}

		//serialize
		msgBytes, err := json.Marshal(msg)
		//fmt.Println(string(msgBytes), msg)
		failOnError(err, "Failed to serialize the message")

		correlationID := RandomString(32)

		//publish
		err = ch.Publish(
			"",           // exchange
			RequestQueue, // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: correlationID,
				ReplyTo:       replyQueue.Name,
				Body:          msgBytes,
			})
		failOnError(err, "Failed to send the message")

		//calculate the time
		start_time := time.Now()
		<-msgs //wait to receive one message
		//fmt.Println(<-msgs)
		total_time := time.Since(start_time)
		arr_times.PushBack(total_time.Seconds()) //store the time

		fmt.Println("Publisher[Default]:", msg.Num)
	}

	mean := calculate_mean(arr_times)
	file, err := os.OpenFile("log-meanTime-RabbitMQProducers.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	failOnError(err, "failed creating file")
	_, err = file.WriteString(fmt.Sprintln(mean))
	failOnError(err, "failed writing to file")
}
