package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	port := "localhost:42069"

	udpAddress, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Printf("Unable to create a new UPD address: %s", err)
	}

	udpConnection, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		fmt.Printf("Unable to create a new UDP connection: %s", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Printf(">")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Unable to read from console: %s", err)
			break
		}

		_, connErr := udpConnection.Write([]byte(line))
		if connErr != nil {
			fmt.Printf("Unable to send data to connection: %s", connErr)
			break
		}

	}
}
