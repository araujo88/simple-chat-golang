package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const TIMEOUT = 60 * 60 // timeout in seconds

var (
	conns   = make(map[net.Conn]bool)
	connsMu sync.Mutex
)

func main() {
	li, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Panic(err)
	}

	defer li.Close()

	for {
		conn, err := li.Accept()

		if err != nil {
			log.Panic(err)
		}

		connsMu.Lock()
		conns[conn] = true
		connsMu.Unlock()

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		connsMu.Lock()
		delete(conns, conn)
		connsMu.Unlock()
	}()

	err := conn.SetDeadline(time.Now().Add(TIMEOUT * time.Second))

	if err != nil {
		log.Println("Failed to set deadline:", err)
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		message := fmt.Sprintf("%s - %s", time.Now().Format("01-02-2006 15:04:05"), line)
		fmt.Println(message)

		broadcast(message)

		// Update the deadline for the next read
		err = conn.SetDeadline(time.Now().Add(TIMEOUT * time.Second))
		if err != nil {
			log.Println("Failed to set deadline:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func broadcast(message string) {
	connsMu.Lock()
	defer connsMu.Unlock()

	for conn := range conns {
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			os.Exit(1)
		}
	}
}
