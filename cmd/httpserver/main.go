package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
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
