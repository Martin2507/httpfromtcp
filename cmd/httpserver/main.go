package main

import (
	"crypto/sha256"
	"encoding/hex"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

const html200 = `<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>`
const html400 = `<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>`
const html500 = `<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>`

func main() {
	server, err := server.Serve(port, handler)
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

func handler(w *response.Writer, req *request.Request) {

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}

	switch req.RequestLine.RequestTarget {

	case "/yourproblem":
		{
			header := response.GetDefaultHeaders(len(html400))
			header["content-type"] = "text/html"

			err := w.WriteStatusLine(response.StatusBadRequest)
			if err != nil {
				log.Printf("Error writing Status Line in httpserver/main.go: %s", err)
				return
			}

			err = w.WriteHeaders(header)
			if err != nil {
				log.Printf("Error writing Headers in httpserver/main.go: %s", err)
				return
			}

			_, err = w.WriteBody([]byte(html400))
			if err != nil {
				log.Printf("Error writing Body in httpserver/main.go: %s", err)
				return
			}
		}

	case "/myproblem":
		{
			header := response.GetDefaultHeaders(len(html500))
			header["content-type"] = "text/html"

			err := w.WriteStatusLine(response.StatusInternalServerError)
			if err != nil {
				log.Printf("Error writing Status Line in httpserver/main.go: %s", err)
				return
			}

			err = w.WriteHeaders(header)
			if err != nil {
				log.Printf("Error writing Headers in httpserver/main.go: %s", err)
				return
			}

			_, err = w.WriteBody([]byte(html500))
			if err != nil {
				log.Printf("Error writing Body in httpserver/main.go: %s", err)
				return
			}
		}

	default:
		{
			header := response.GetDefaultHeaders(len(html200))
			header["content-type"] = "text/html"

			err := w.WriteStatusLine(response.StatusOK)
			if err != nil {
				log.Printf("Error writing Status Line in httpserver/main.go: %s", err)
				return
			}

			err = w.WriteHeaders(header)
			if err != nil {
				log.Printf("Error writing Headers in httpserver/main.go: %s", err)
				return
			}

			_, err = w.WriteBody([]byte(html200))
			if err != nil {
				log.Printf("Error writing Body in httpserver/main.go: %s", err)
				return
			}
		}
	}

}

func proxyHandler(w *response.Writer, req *request.Request) {

	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	fullPath := "https://httpbin.org/" + path

	resp, err := http.Get(fullPath)
	if err != nil {
		log.Printf("Error: Unable to connect to httpbin.org: %s", err)
		return
	}

	defer resp.Body.Close()

	header := headers.NewHeaders()
	header["connection"] = "close"
	header["content-type"] = resp.Header.Get("Content-Type")
	header["Transfer-Encoding"] = "chunked"
	header["Trailer"] = "X-Content-SHA256, X-Content-Length"

	err = w.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Printf("Error writing Status Line in httpserver/main.go: %s", err)
		return
	}

	err = w.WriteHeaders(header)
	if err != nil {
		log.Printf("Error writing Headers in httpserver/main.go: %s", err)
		return
	}

	buffer := make([]byte, 1024)
	var fullBody []byte

	for {

		n, err := resp.Body.Read(buffer)

		if n > 0 {
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				log.Printf("Error writing Body in httpserver/main.go: %s", err)
				return
			}

			fullBody = append(fullBody, buffer[:n]...)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error: Unable to read data from response: %s", err)
			break
		}

	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("Error: Unable to update the state to done: %s", err)
		return
	}

	trailerHeaders := headers.NewHeaders()

	hashHex := sha256.Sum256(fullBody)
	hashStr := hex.EncodeToString(hashHex[:])

	trailerHeaders["X-Content-SHA256"] = hashStr
	trailerHeaders["X-Content-Length"] = strconv.Itoa(len(fullBody))

	err = w.WriteTrailers(trailerHeaders)
	if err != nil {
		log.Printf("Error: Unable to write a trailer: %s", err)
		return
	}
}
