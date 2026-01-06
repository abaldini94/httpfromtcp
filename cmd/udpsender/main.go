package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":42069")

	if err != nil {
		fmt.Println("Failed to resolve UDP addr")
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Failed to create udp connection")
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		newLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading a new line from Stdin")
		}

		_, err = udpConn.Write([]byte(newLine))
		if err != nil {
			fmt.Printf("Error while writing: %s", err.Error())
		}
	}
}
