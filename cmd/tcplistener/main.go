package main

import (
	"fmt"
	"log"
	"net"

	"github.com/JA50N14/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %v\n", err)
	}

	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %v\n", err)
		}
		fmt.Println("Connection has been accepted from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error parsing request: %s\n", err.Error())
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
	}
}

// func getLinesChannel(conn io.ReadCloser) <-chan string {
// 	lines := make(chan string)

// 	go func() {
// 		defer conn.Close()
// 		defer close(lines)
// 		currentLine := ""
// 		for {
// 			buff := make([]byte, 8)
// 			n, err := conn.Read(buff)
// 			if err != nil {
// 				if currentLine != "" {
// 					lines <- currentLine
// 					currentLine = ""
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break
// 				}
// 				fmt.Printf("error: %s\n", err.Error())
// 				break
// 			}
// 			str := string(buff[:n])
// 			parts := strings.Split(str, "\n")

// 			for i := 0; i < len(parts)-1; i++ {
// 				lines <- currentLine + parts[i]
// 				currentLine = ""
// 			}
// 			currentLine += parts[len(parts)-1]
// 		}
// 	}()

// 	return lines
// }
