package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {

	// Open the input file for streaming line-by-line reads.
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("unable to open %s: %s\n", inputFilePath, err)
	}

	// getLinesChannel starts a goroutine that reads from the file,
	// emits complete lines, then closes both the file and channel.
	lines := getLinesChannel(file)

	// Consume each emitted line and print it in the required format.
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		// The reader owns cleanup because it may still be reading
		// after this function returns.
		defer file.Close()
		defer close(lines)

		currentLineContent := ""

		for {
			// Read fixed-size chunks so lines may arrive in pieces.
			buffer := make([]byte, 8)
			n, err := file.Read(buffer)

			if err != nil {
				// If EOF arrives after partial line content, emit it first.
				if currentLineContent != "" {
					lines <- currentLineContent
				}

				if errors.Is(err, io.EOF) {
					break
				}

				fmt.Printf("unable to read file: %v\n", err)
				return
			}

			// Only convert the bytes that were actually read.
			str := string(buffer[:n])

			// Split on newline. All parts except the last are complete lines.
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- currentLineContent + parts[i]
				currentLineContent = ""
			}

			// The final part may be an unfinished line, so keep buffering it.
			currentLineContent += parts[len(parts)-1]
		}
	}()

	return lines
}
