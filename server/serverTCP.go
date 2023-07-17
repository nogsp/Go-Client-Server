package main

import (
	//"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

func main() {

	r, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	listen, err:= net.ListenTCP("tcp", r)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer listen.Close()

	println("Servidor pronto para receber mensagens TCP.")


	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)

	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			println(err)
			break // Exit the loop when encountering an error
		}
		n, err := strconv.Atoi(string(buffer[:bytesRead]))

		//fmt.Printf("Recebido de %s: %d\n", conn.LocalAddr().String(), n)
		// write data to response
		_, err = conn.Write([]byte(strconv.Itoa(fibo(n))))
		if err != nil {
			println(err)
			break // Exit the loop when encountering an error
		}
	}

	conn.Close() // Close the connection after handling the request
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
