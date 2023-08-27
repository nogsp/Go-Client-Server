package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const MQTTHost = "tcp://172.17.0.3:1883"
const MQTTTopic = "PubSub"

type Message struct {
	Msg string `json:"msg"`
	Pid int    `json:"pid"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func fibo(n int) int {
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

var callback MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	num, err := strconv.Atoi(string(m.Payload()))
	failOnError(err, "Failed to convert string to int")

	fmt.Println("Received Message:", m.Payload(), "\tTopic:", m.Topic())

	var msg Message
	err = json.Unmarshal(m.Payload(), &msg)
	failOnError(err, "Failed to deserialize")
	msg.Msg = strconv.Itoa(fibo(num))

	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID("publisher")

	client := MQTT.NewClient(opts)
	token := client.Connect()
	token.Wait()
	failOnError(token.Error(), "Failed to connect to the broker")

	payload, err := json.Marshal(msg)
	failOnError(err, "Failed to marshal message")

	client.Publish(
		MQTTTopic+"/"+strconv.Itoa(msg.Pid),
		2,
		false,
		payload,
	)
	client.Disconnect(250)
}

func main() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID("subscriber 0")

	client := MQTT.NewClient(opts)
	token := client.Connect()
	token.Wait()
	failOnError(token.Error(), "Failed to connect to the broker")

	defer client.Disconnect(250)

	token = client.Subscribe(
		MQTTTopic,
		2,
		callback,
	)
	token.Wait()
	failOnError(token.Error(), "Failed to subscribe to topic")

	fmt.Println("Consumer starts")

	// Wait indefinitely
	select {}
}
