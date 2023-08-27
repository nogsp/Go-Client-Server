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

const MQTTHost = "tcp://0.0.0.0:1883"
const MQTTTopic = "PubSub"

const SampleSize = 10000

type Message struct {
	Msg string `json:"msg"`
	Pid int    `json:"pid"`
}

var (
	start      time.Time
	total_time time.Duration
	arr_times  *list.List
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func subs(PID int) {
	// configure client
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID("subscriber " + strconv.Itoa(PID))

	// create client
	client := MQTT.NewClient(opts)

	// connect to the broker
	token := client.Connect()
	token.Wait()
	failOnError(token.Error(), "Failed to connect to the broker")

	defer client.Disconnect(250)

	callback := func(c MQTT.Client, m MQTT.Message) {
		fmt.Println("Received Message:", m.Payload(), "\tTopic:", m.Topic())

		total_time = time.Since(start)
		arr_times.PushBack(total_time.Seconds())

		var msg Message
		err := json.Unmarshal(m.Payload(), &msg)
		failOnError(err, "Failed to deserialize")
		fmt.Println(msg.Msg)
	}

	token = client.Subscribe(
		MQTTTopic+"/"+strconv.Itoa(os.Getpid()),
		2,
		callback,
	)
	token.Wait()
	failOnError(token.Error(), "Failed to subscribe to topic")
}

func calculate_mean(l *list.List) float64 {
	total := 0.0
	for e := l.Front(); e != nil; e = e.Next() {
		if num, ok := e.Value.(float64); ok {
			total += num
		}
	}
	return total / float64(l.Len())
}

func main() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTHost)
	opts.SetClientID("publisher")

	client := MQTT.NewClient(opts)

	token := client.Connect()
	token.Wait()
	failOnError(token.Error(), "Failed to connect to the broker")

	PID := os.Getpid()

	arr_times = list.New()

	subs(PID)

	defer client.Disconnect(250)

	for i := 0; i < SampleSize; i++ {
		msg := Message{
			Msg: fmt.Sprint(i % 50),
			Pid: PID,
		}

		payload, err := json.Marshal(msg)
		failOnError(err, "Failed to marshal message")

		client.Publish(
			MQTTTopic,
			2,
			false,
			payload,
		)

		start = time.Now()

		fmt.Println("Published Message:", msg)

		time.Sleep(time.Microsecond)
	}

	mean := calculate_mean(arr_times)
	file, err := os.OpenFile("log-meanTime-MQTTPublishers.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	failOnError(err, "Failed to create file")
	_, err = file.WriteString(fmt.Sprintln(mean))
	failOnError(err, "Failed to write to file")
}
