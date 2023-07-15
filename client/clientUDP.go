package main

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8082"
	TYPE = "udp"
)

func main() {
	servAddr, err := net.ResolveUDPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("Erro ao resolver endereço do servidor:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, servAddr)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	arr_times := list.New()
	total_iterations := 10000
	aplication_time := time.Now()
	for iteration := 0; iteration < total_iterations; iteration++ {
		num := strconv.Itoa(rand.Intn(80))
		//fmt.Printf("Fibo for %s\n", num)
		start_time := time.Now()
		_, err = conn.Write([]byte(num))
		if err != nil {
			fmt.Println("Erro ao enviar mensagem:", err)
			os.Exit(1)
		}

		received := make([]byte, 1024)
		_, _, err := conn.ReadFromUDP(received)
		if err != nil {
			fmt.Println("Erro ao receber resposta:", err)
			os.Exit(1)
		}
		total_time := time.Since(start_time)
		arr_times.PushBack(total_time.Seconds())

		//fmt.Printf("Fibo for %s is %s\n", num, string(received))
		//fmt.Printf("Took %f second\n", total_time.Seconds())

	}
	conn.Close()
	println("Tempo total:", time.Since(aplication_time).Seconds())
	println("Tempo medio:", calculate_mean(arr_times))
	println("Desvio padrao:", calculate_deviation(arr_times))
	println()
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
