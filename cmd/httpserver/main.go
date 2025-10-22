package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JA50N14/httpfromtcp/internal/headers"
	"github.com/JA50N14/httpfromtcp/internal/request"
	"github.com/JA50N14/httpfromtcp/internal/response"
	"github.com/JA50N14/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	h := headers.NewHeaders()
	h.Set("Content-Type", "text/html")

	responseBody := getResponseBody(req.RequestLine.RequestTarget)

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.StatusCode = response.StatusCodeBadRequest
	case "/myproblem":
		w.StatusCode = response.StatusCodeInternalServerError
	case "/":
		w.StatusCode = response.StatusCodeSuccess
	}

	w.BodyLen = len(responseBody)

	w.WriteStatusLine(w.StatusCode)
	w.WriteHeaders(h)
	w.WriteBody([]byte(responseBody))

	//TESTING///////////////////////////////////////////
	fmt.Println("REQEST/RESPONSE VALUES")
	fmt.Printf("StatusCode: %d\n", w.StatusCode)
	fmt.Printf("Headers:\n")
	for key, value := range w.Headers {
		fmt.Printf("Key: %s / Value: %s\n", key, value)
	}
	fmt.Println("Body:")
	fmt.Printf("%s", string(w.Body))
	fmt.Println()
	fmt.Println()
}

func getResponseBody(requestTarget string) string {
	if requestTarget == "/yourproblem" {
		return `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	}

	if requestTarget == "/myproblem" {
		return `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	}

	return `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
}

// func handler(w io.Writer, req *request.Request) *server.HandlerError {
// 	if req.RequestLine.RequestTarget == "/yourproblem" {
// 		return &server.HandlerError{
// 			StatusCode: response.StatusCodeBadRequest,
// 			Message:    "Your problem is not my problem\n",
// 		}
// 	}
// 	if req.RequestLine.RequestTarget == "/myproblem" {
// 		return &server.HandlerError{
// 			StatusCode: response.StatusCodeInternalServerError,
// 			Message:    "Woopsie, my bad\n",
// 		}
// 	}
// 	w.Write([]byte("All good, frfr\n"))
// 	return nil
// }
