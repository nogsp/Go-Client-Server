package main

import (
	"container/list"
	"fmt"
	"log"

	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	arr_times := list.New()
	total_iterations := 10000
	aplication_time := time.Now()

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	for iteration := 0; iteration < total_iterations; iteration++ {
		num := strconv.Itoa(rand.Intn(80))

		start_time := time.Now()
		_, err = conn.Write([]byte(num))
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}

		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}

		total_time := time.Since(start_time)
		arr_times.PushBack(total_time.Seconds())

		//fmt.Printf("Fibo for %s is %s\n", num, string(received))
		//fmt.Printf("Took %f second\n", total_time.Seconds())
		//time.Sleep(time.Second)
	}
	conn.Close()
	println("Tempo total:", time.Since(aplication_time).Seconds())
	mean := calculate_mean(arr_times)
	println("Tempo medio:", mean)
	println("Desvio padrao:", calculate_deviation(arr_times))
	println()
	file, err := os.OpenFile("log-meanTime-TCPClients.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	_, err = file.WriteString(fmt.Sprintln(mean))
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	//println("Desvio padrao:", calculate_deviation(&arr_times))
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

func calculate_deviation(l *list.List) float64 {
	mean := calculate_mean(l)
	deviation := 0.0

	for e := l.Front(); e != nil; e = e.Next() {
		if num, ok := e.Value.(float64); ok {
			deviation += math.Pow(num-mean, 2)
		}
	}
	n := float64(l.Len())
	return math.Pow(deviation/n, 0.5)
}
