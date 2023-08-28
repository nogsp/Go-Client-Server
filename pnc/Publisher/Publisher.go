package main

import (
	"encoding/json"
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const MQTTHost = "tcp://0.0.0.0:1883"
const MQTTTopic = "Fibonacci"

type Message struct {
	Msg string `json:"msg"`
	Pid int    `json:"pid"`
}

var receiveHandler MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	fmt.Printf("Mensagem recebida, o fibo eh %s\n", m.Payload())
}

func main() {
	// config
	clientID := "publisher_" + string(os.Getpid())
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID(clientID)

	// criar cliente
	client := MQTT.NewClient(opts)

	// conectar ao broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Subscreve para a resposta
	if token := client.Subscribe(MQTTTopic+"/"+clientID, 2, receiveHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		msg := Message{
			Msg: fmt.Sprint((i + 1) % 50),
			Pid: os.Getpid(),
		}
		jmsg, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		// Publicar a mensagem
		token := client.Publish(MQTTTopic, 2, false, jmsg)
		token.Wait()
		if token.Error() != nil {
			panic(token.Error())
		}

		fmt.Printf("Mensagem publicada %s com PID %d\n", msg.Msg, msg.Pid)
		//time.Sleep(time.Second)
	}

	// Desconectar do Broker no final
	client.Disconnect(250)

}
