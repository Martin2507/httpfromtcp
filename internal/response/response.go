package response

import (
	"errors"
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
