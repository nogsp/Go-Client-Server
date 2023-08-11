package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

	conn, err := net.ListenUDP(TYPE, servAddr)
	if err != nil {
		fmt.Println("Erro ao ouvir:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Servidor pronto para receber mensagens UDP.")

	//buffer := make([]byte, 1024)

	for {
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao receber dados:", err)
			os.Exit(1)
		}
		go handleRequest(*conn, n, addr, buffer)
	}
}

func handleRequest(conn net.UDPConn, n int, addr *net.UDPAddr, buffer []byte) {
	num, err := strconv.Atoi(string(buffer[:n]))
	//fmt.Printf("Recebido de %s: %d\n", addr.String(), num)

	// Enviar uma resposta ao cliente
	response := []byte(strconv.Itoa(fibo(num)))
	_, err = conn.WriteToUDP(response, addr)
	if err != nil {
		fmt.Println("Erro ao enviar resposta:", err)
		os.Exit(1)
	}
}

func fibo(n int) int {
	ans := 1
	prev := 0
	for i := 1; i < n; i++ {
		temp := ans
		ans = ans + prev
		prev = temp
	}
	return ans
}
