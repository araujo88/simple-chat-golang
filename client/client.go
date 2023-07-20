package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var name string

func sendMessage(conn net.Conn) {
	//fmt.Print("Enter text: \n")
	reader := bufio.NewReader(os.Stdin)
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSuffix(msg, "\n") // Removing the trailing newline
	_, err := conn.Write([]byte(name + ": " + msg + "\n"))
	if err != nil {
		fmt.Println("Error sending message:", err)
		os.Exit(1)
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}

func main() {
	fmt.Print("Enter your name: ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Create a connection
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	go handleConnection(conn)

	for {
		sendMessage(conn)
	}
}
