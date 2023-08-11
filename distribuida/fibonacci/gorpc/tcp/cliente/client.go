package main

import (
	"aulas/distribuida/shared"
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

func main() {

	var reply int

	// conecta ao consumer (Fibonacci)
	client, err := rpc.Dial("tcp", ":"+strconv.Itoa(shared.FibonacciPort))
	shared.ChecaErro(err, "Não foi possível criar uma conexão TCP para o consumer Fibonacci...")

	defer func(client *rpc.Client) {
		var err = client.Close()
		shared.ChecaErro(err, "Não foi possível fechar a conexão TCP com o consumer Fibonacci...")
	}(client)

	// invoca operação remota do Fibonacci

	total_iterations := 10000
	arr_times := list.New()

	for iteration := 0; iteration < total_iterations; iteration++ {
		num := rand.Intn(80)
		start_time := time.Now()
		err = client.Call("Fibonacci.Fibo", num, &reply)
		shared.ChecaErro(err, "Erro na invocação remota do Fibonacci...")
		total_time := time.Since(start_time)
		arr_times.PushBack(total_time.Seconds())
		fmt.Printf("Fibo(%v) = %v \n", num, reply)
	}

	mean := calculate_mean(arr_times)

	file, err := os.OpenFile("log-meanTime-RPCClients.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	_, err = file.WriteString(fmt.Sprintln(mean))
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
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
