package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {

	port := ":42069"

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Unable to create a new listener: %s", err)
	}

	defer listener.Close()

	for {

		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Unable to create a new connection: %s\n", err)
			break
		}

		fmt.Println("A new connection has been created and accepted")

		// // getLinesChannel starts a goroutine that reads from the file, emits complete lines, then closes both the file and channel.
		// lines := getLinesChannel(connection)

		// // Consume each emitted line and print it in the required format.
		// for line := range lines {
		// 	fmt.Println(line)
		// }

		req, err := request.RequestFromReader(connection)

		if err != nil {
			fmt.Printf("Unable to read data from provided source")
			break
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		fmt.Println("Connection has been closed")
	}

}

// func getLinesChannel(file io.ReadCloser) <-chan string {
// 	lines := make(chan string)

// 	go func() {
// 		// The reader owns cleanup because it may still be reading after this function returns.
// 		defer file.Close()
// 		defer close(lines)

// 		currentLineContent := ""

// 		for {
// 			// Read fixed-size chunks so lines may arrive in pieces.
// 			buffer := make([]byte, 8)
// 			n, err := file.Read(buffer)

// 			if err != nil {
// 				// If EOF arrives after partial line content, emit it first.
// 				if currentLineContent != "" {
// 					lines <- currentLineContent
// 				}

// 				if errors.Is(err, io.EOF) {
// 					break
// 				}

// 				fmt.Printf("unable to read file: %v\n", err)
// 				return
// 			}

// 			// Only convert the bytes that were actually read.
// 			str := string(buffer[:n])

// 			// Split on newline. All parts except the last are complete lines.
// 			parts := strings.Split(str, "\n")

// 			for i := 0; i < len(parts)-1; i++ {
// 				lines <- currentLineContent + parts[i]
// 				currentLineContent = ""
// 			}

// 			// The final part may be an unfinished line, so keep buffering it.
// 			currentLineContent += parts[len(parts)-1]
// 		}
// 	}()

// 	return lines
// }
