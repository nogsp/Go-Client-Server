package main

import (
	"log"
	"net"
	"os"
	"strconv"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}

}
func add(a int, b int) int {
	return a + b
}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	n, err := strconv.Atoi(string(buffer[:bytesRead]))
	ans := 1
	prev := 0
	for i := 1; i < n; i++ {
		temp := ans
		ans = ans + prev
		prev = temp
	}

	//fmt.Println(ans)
	// write data to response
	conn.Write([]byte(strconv.Itoa(ans)))

	// close conn
	conn.Close()
}
