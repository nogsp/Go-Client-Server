package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	for {

		n := strconv.Itoa(rand.Intn(80))
		fmt.Printf("Fibo for %s\n", n)
		conn, err := net.DialTCP(TYPE, nil, tcpServer)
		if err != nil {
			println("Dial failed:", err.Error())
			os.Exit(1)
		}

		_, err = conn.Write([]byte(n))
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}

		// buffer to get data
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}

		println("Received message:", string(received))

		conn.Close()
		time.Sleep(time.Second)
	}
}
