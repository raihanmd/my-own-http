package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udp, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	bufioReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := bufioReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte(line))
	}
}
