package main

import (
	"fmt"
	"net"

	"httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Impossible to read the file")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Impossible to accept connection")
		}

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Impossible to get request from reader")
		}

		outStr := "Request line:\n"
		outStr += fmt.Sprintf("- Method: %s\n", req.RequestLine.Method)
		outStr += fmt.Sprintf("- Target: %s\n", req.RequestLine.RequestTarget)
		outStr += fmt.Sprintf("- Version: %s\n", req.RequestLine.HttpVersion)
		outStr += "Headers:\n"
		for name, value := range req.Headers {
			outStr += fmt.Sprintf("%s: %s\n", name, value)
		}
		outStr += fmt.Sprintf("Body:\n%s\n", string(req.Body))
		fmt.Println(outStr)
	}
}
