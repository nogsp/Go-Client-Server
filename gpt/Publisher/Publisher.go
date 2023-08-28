package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("Message %d", i+1)
		token := client.Publish("topic", 0, false, message)
		token.Wait()
		log.Printf("Sent: %s\n", message)
		time.Sleep(time.Second)
	}

	client.Disconnect(250)
}
