package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MQTT
const MQTTHost = "mqtt://172.17.0.3:1883"
const MQTTTopic = "PubSub"
const MQTTRequest = "request"
const MQTTReply = "reply"

// Other configurations
const SampleSize = 10000

var (
	start      time.Time
	total_time time.Duration
	arr_times  *list.List
)

type Message struct {
	Msg string
	Pid int
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func subs(PID int) {
	//configure client
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID("subscriber " + strconv.Itoa(PID))

	//Create client
	client := MQTT.NewClient(opts)

	//connect to the broker
	token := client.Connect()
	token.Wait()
	failOnError(token.Error(), "Failed to connect to the broker")

	//Disconnect to the broker
	defer client.Disconnect(250)

	callback := func(c MQTT.Client, m MQTT.Message) {
		//process the message
		fmt.Println("Mensagem Recebida:", m.Payload(), "\tTopic:", m.Topic())

		//calculate the time
		total_time = time.Since(start)
		arr_times = list.New()

		arr_times.PushBack(total_time.Seconds()) //store the time

		var msg Message
		err := json.Unmarshal(m.Payload(), &msg)
		failOnError(err, "It was not possible to do deserialization")
		fmt.Println(msg.Msg)

	}

	token = client.Subscribe(
		MQTTTopic+"/"+strconv.Itoa(os.Getpid()), //Topic
		2,        //Quality of Service
		callback, //callback
	)
	token.Wait()
	failOnError(token.Error(), "Failed to subscribe in topic")
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
	//configure client
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID("publisher")

	//Create client
	client := MQTT.NewClient(opts)

	//connect to the broker
	token := client.Connect()
	token.Wait()
	failOnError(token.Error(), "Failed to connect to the broker")

	//Get PID
	PID := os.Getpid()

	//Subscribe to receive responses
	subs(PID)

	//Disconnect to the broker
	defer client.Disconnect(250)

	for i := 0; i < SampleSize; i++ {

		//create the message
		msg := Message{
			Msg: fmt.Sprint(i),
			Pid: PID,
		}

		//publish the message
		client.Publish(
			MQTTTopic, //Topic
			2,         //Quality of Service
			false,     //retained
			msg,       //payload
		)

		start = time.Now()

		fmt.Println("Mensagem publicada:", msg)

		time.Sleep(time.Microsecond)
	}

	mean := calculate_mean(arr_times)
	file, err := os.OpenFile("log-meanTime-MQTTPublishers.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	failOnError(err, "failed creating file")
	_, err = file.WriteString(fmt.Sprintln(mean))
	failOnError(err, "failed writing to file")
}
