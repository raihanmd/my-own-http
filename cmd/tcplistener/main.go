package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)

		data := make([]byte, 8)
		currentLine := ""

		for {
			count, err := f.Read(data)
			if err != nil {
				if err == io.EOF {
					if currentLine != "" {
						lines <- currentLine
					}
					break
				}
				return
			}

			currentLine += string(data[:count])
			parts := strings.Split(currentLine, "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- parts[i]
			}

			currentLine = parts[len(parts)-1]
		}
	}()

	return lines
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	// fmt.Println("Listening on :42069...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// fmt.Println("Accepted connection")
		go func(c net.Conn) {
			for line := range getLinesChannel(c) {
				fmt.Println(line)
			}
			// fmt.Println("Connection closed")
		}(conn)
	}
}
