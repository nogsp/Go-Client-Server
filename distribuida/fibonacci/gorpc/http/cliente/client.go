package main

import (
	"aulas/distribuida/shared"
	"fmt"
	"net/rpc"
	"strconv"
)

func Cliente() {
	var reply int

	// conecta ao consumer (Fibonacci)
	client, err := rpc.DialHTTP("tcp", ":"+strconv.Itoa(shared.FibonacciPort))
	shared.ChecaErro(err, "Não foi possível criar uma conexão com o consumer Fibonacci...")

	defer func(c *rpc.Client) {
		var err = c.Close()
		shared.ChecaErro(err, "Erro ao fechar a conexão com o consumer Fibonacci...")
	}(client)

	// invoca operação remota do fibonacci
	n := 10
	err = client.Call("Fibonacci.Fibo", n, &reply)
	shared.ChecaErro(err, "Erro na invocação remota do consumer Fibonacci...")
	fmt.Printf("Fibo(%v) = %v \n", n, reply)
}

func main() {

	go Cliente()

	fmt.Scanln()
}
