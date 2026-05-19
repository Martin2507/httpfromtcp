package response

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type responseState int

const (
	responseStateWriteStatusLine responseState = iota
	responseStateWriteHeaders
	responseStateWriteBody
	responseStateWriteChunkBody
	responseStateDone
)

type Writer struct {
	W     io.Writer
	state responseState
}

func GetDefaultHeaders(contentLen int) headers.Headers {

	h := headers.NewHeaders()

	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"

	return h
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state == responseStateWriteStatusLine {
		switch statusCode {

		case StatusOK:
			{
				_, err := w.W.Write([]byte("HTTP/1.1 200 OK\r\n"))
				if err != nil {
					return err
				}
			}

		case StatusBadRequest:
			{
				_, err := w.W.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
				if err != nil {
					return err
				}
			}

		case StatusInternalServerError:
			{
				_, err := w.W.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
				if err != nil {
					return err
				}
			}

		default:
			{
				formattedString := "HTTP/1.1 " + strconv.Itoa(int(statusCode)) + " \r\n"
				_, err := w.W.Write([]byte(formattedString))
				if err != nil {
					return err
				}
			}

		}

		w.state = responseStateWriteHeaders

		return nil
	}

	return errors.New("Error: Action out of order")
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.state == responseStateWriteHeaders {
		for key, val := range headers {

			formattedString := key + ": " + val + "\r\n"

			_, err := w.W.Write([]byte(formattedString))

			if err != nil {
				return err
			}

		}

		_, err := w.W.Write([]byte("\r\n"))

		if err != nil {
			return err
		}

		w.state = responseStateWriteBody

		return nil
	}

	return errors.New("Error: Action out of order")
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.state == responseStateWriteBody {

		w.state = responseStateDone

		return w.W.Write(p)
	}

	return 0, errors.New("Error: Action out of order")

}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {

	return fmt.Fprintf(w.W, "%x\r\n%s\r\n", len(p), p)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {

	if w.state == responseStateDone {
		return 0, errors.New("Error: Unable to perform action state is already Done")
	}

	w.state = responseStateWriteChunkBody

	return fmt.Fprint(w.W, "0\r\n")

}

func (w *Writer) WriteTrailers(h headers.Headers) error {

	if w.state != responseStateWriteChunkBody {
		return errors.New("Error: Unable to perform actions out of order, writing trailer must be the last action")
	}

	for key, val := range h {

		formattedHeader := key + ": " + val + "\r\n"

		_, err := w.W.Write([]byte(formattedHeader))

		if err != nil {
			return err
		}
	}

	_, err := w.W.Write([]byte("\r\n"))

	if err != nil {
		return err
	}

	w.state = responseStateDone

	return nil

}
