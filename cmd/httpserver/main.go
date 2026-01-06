package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func respond200() string {
	return `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Everything OK</p>
  </body>
</html>`
}

func respond400() string {
	return `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Bad, Bad, Bad</p>
  </body>
</html>`

}

func respond500() string {
	return `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>My bad, My bad, My bad</p>
  </body>
</html>`

}

func handler(w *response.Writer, r *request.Request) {
	var body string
	target := r.RequestLine.RequestTarget
	h := response.GetDefaultHeaders(0)
	switch target {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusBadRequest)
		body = respond400()
	case "/myproblem":
		w.WriteStatusLine(response.StatusServerError)
		body = respond500()
	default:
		w.WriteStatusLine(response.StatusOK)
		body = respond200()
	}
	b := []byte(body)
	h.Replace("content-length", fmt.Sprintf("%d", len(b)))
	h.Replace("content-type", "html")
	w.WriteHeaders(h)
	w.WriteBody(b)
}

func main() {
	server, err := Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
