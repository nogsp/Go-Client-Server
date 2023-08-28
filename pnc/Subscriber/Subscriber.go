package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const MQTTHost = "mqtt://172.17.0.5:1883"
const MQTTTopic = "Fibonacci"

type Message struct {
	Msg string `json:"msg"`
	Pid int    `json:"pid"`
}

func main() {
	// Configurar cliente
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID("subscriber 1")

	// Criar cliente
	client := MQTT.NewClient(opts)

	// Conectar ao Broker
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		panic(token.Error())
	}

	//Subscreve ao tópico e definir handler
	token = client.Subscribe(MQTTTopic, 2, func(c MQTT.Client, m MQTT.Message) {
		var msg Message
		json.Unmarshal(m.Payload(), &msg)
		ans := fmt.Sprintf("Mensagem recebida, o fibo de %s eh %s\n", msg.Msg, fibo(msg.Msg))
		fmt.Printf("Mensagem recebida, o fibo de %s eh %s\n", msg.Msg, fibo(msg.Msg))

		token := client.Publish(MQTTTopic+"/publisher_"+fmt.Sprint(msg.Pid), 2, false, ans)
		token.Wait()
		if token.Error() != nil {
			panic(token.Error())
		}
	})
	token.Wait()
	if token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Consumidor on")
	fmt.Scanln()

	// Desconectar  do brker no final
	client.Disconnect(250)

}

func fibo(strNum string) string {
	n, err := strconv.Atoi(strNum)
	if err != nil {
		panic(err)
	}
	ans := 1
	prev := 0
	for i := 1; i < n; i++ {
		temp := ans
		ans = ans + prev
		prev = temp
	}
	return strconv.Itoa(ans)
}
