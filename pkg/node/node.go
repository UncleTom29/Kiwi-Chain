package main

import (
	"sync"
	"time"
	"net"
	"log"
	"io"
	"bufio"
)



func handleConnection(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Welcome to Kiwi Blockchain!\n")

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := scanner.Text()
		switch msg {
		case "get blockchain":
			broadcastChain()
		default:
			io.WriteString(conn, "\nEnter a new BPM:")

			bpm := make([]byte, 10)
			_, err := conn.Read(bpm)
			if err != nil {
				log.Println(err)
				return
			}

			newBlock, err := createBlock(Blockchain[len(Blockchain)-1], []Transaction{}, "validatorAddress")
			if err != nil {
				log.Println(err)
				return
			}

			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				candidateBlocks <- newBlock
			}

			io.WriteString(conn, fmt.Sprintf("\n%s", Blockchain))
		}
	}
}
