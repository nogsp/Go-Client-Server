package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const MQTTHost = "mqtt://172.17.0.5:1883"
const MQTTTopic = "Fibonacci"
const QoS = 0

type Message struct {
	Msg string `json:"msg"`
	Pid int    `json:"pid"`
}

var m runtime.MemStats

var start time.Time
var total_time time.Duration

var signal bool

var receiveHandler MQTT.MessageHandler = func(c MQTT.Client, m MQTT.Message) {
	total_time = time.Since(start)
	fmt.Println(total_time)
	signal = true
	fmt.Printf("Mensagem recebida, o fibo eh %s\n", m.Payload())
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

func getCPUSample() (idle, total uint64) {
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func main() {
	// config
	clientID := "publisher_" + fmt.Sprint(os.Getpid())
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
	if token := client.Subscribe(MQTTTopic+"/"+clientID, QoS, receiveHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	arr_times := list.New()
	arr_mems := list.New()

	//get initial CPU informations
	idle0, total0 := getCPUSample()

	for i := 0; i < 10000; i++ {
		signal = false
		msg := Message{
			Msg: fmt.Sprint((i + 1) % 50),
			Pid: os.Getpid(),
		}
		jmsg, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		// Publicar a mensagem
		token := client.Publish(MQTTTopic, QoS, false, jmsg)
		start = time.Now()
		token.Wait()
		if token.Error() != nil {
			panic(token.Error())
		}

		// Esperar resposta
		for signal == false {

		}
		arr_times.PushBack(total_time.Seconds())

		//get the Memory obtained from Sys(MiB)
		runtime.ReadMemStats(&m)
		arr_mems.PushBack(float64(m.Sys) / 1024 / 1024)

		fmt.Printf("Mensagem publicada %s com PID %d\n", msg.Msg, msg.Pid)
		//time.Sleep(time.Second)
	}

	// Desconectar do Broker no final
	client.Disconnect(250)

	//get the final CPU informations
	idle1, total1 := getCPUSample()
	//to calculate CPU Usage
	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	//save the cpu Usage
	file, err := os.OpenFile("log-CPUUsage-MQTTPublishers.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(fmt.Sprintln(cpuUsage))
	if err != nil {
		panic(err)
	}

	//save the mean time
	mean := calculate_mean(arr_times)
	file, err = os.OpenFile("log-meanTime-MQTTPublishers.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(fmt.Sprintln(mean))
	if err != nil {
		panic(err)
	}

	//save the mean memory obtained from OS
	mean = calculate_mean(arr_mems)
	file, err = os.OpenFile("log-meanMemSys-MQTTPublishers.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(fmt.Sprintln(mean))
	if err != nil {
		panic(err)
	}

}
