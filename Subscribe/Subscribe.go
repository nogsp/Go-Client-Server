package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	opts := mqtt.NewClientOptions().AddBroker("tcp://0.0.0.0:1883")
	opts.SetClientID("subscriber")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Subscribe("topic", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})

	token.Wait()

	<-c
}
