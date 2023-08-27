package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MQTT
const MQTTHost = "mqtt://172.17.0.3:1883"
const MQTTTopic = "PubSub"
const MQTTRequest = "request"
const MQTTReply = "reply"

// Other configurations
const SampleSize = 10000

type Message struct {
	Msg string
	Pid int
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

var callback MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	//process the message
	num, err := strconv.Atoi(string(m.Payload()))
	failOnError(err, "It was not possible to convert string to int")
	fmt.Println("Mensagem Recebida:", m.Payload(), "\tTopic:", m.Topic())
	var msg Message
	err = json.Unmarshal(m.Payload(), &msg)
	failOnError(err, "It was not possible to do deserialization")
	msg.Msg = strconv.Itoa(fibo(num))

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

	client.Publish(
		MQTTTopic+"/"+strconv.Itoa(msg.Pid), //Topic
		2,                                   //Quality of Service
		false,                               //retained
		msg,                                 //payload
	)
	client.Disconnect(250)
}

func main() {
	//configure client
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID("subscriber 0")

	//Create client
	client := MQTT.NewClient(opts)

	//connect to the broker
	token := client.Connect()
	token.Wait()
	failOnError(token.Error(), "Failed to connect to the broker")

	//Disconnect to the broker
	defer client.Disconnect(250)

	token = client.Subscribe(
		MQTTTopic, //Topic
		2,         //Quality of Service
		callback,  //callback
	)
	token.Wait()
	failOnError(token.Error(), "Failed to subscribe in topic")

	fmt.Println("Consumer starts")
}
