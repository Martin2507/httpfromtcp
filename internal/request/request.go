package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestInitialized requestState = iota
	requestDone
)

const bufferSize int = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	// data, err := io.ReadAll(reader)
	// if err != nil {
	// 	fmt.Printf("Unexpected error has occured while reading data: %s", err)
	// 	return nil, err
	// }

	// result, _, err := parseRequestLine(data)
	// if err != nil {
	// 	fmt.Printf("Unexpected error has occured while parsing data: %s", err)
	// 	return &Request{}, err
	// }

	// return result, nil

	buffer := make([]byte, bufferSize)

	readToIndex := 0

	req := Request{}

	req.state = 0

	for req.state != 1 {

		if readToIndex == len(buffer) {

			temp := make([]byte, len(buffer)*2)
			copy(temp, buffer)

			buffer = temp
		}

		readBytes, err := reader.Read(buffer[readToIndex:])

		if errors.Is(err, io.EOF) {
			req.state = 1
			break
		} else if err != nil {
			fmt.Printf("Unexpected error has occured while reading data: %s", err)
			return &Request{}, err
		}

		readToIndex += readBytes

		count, err := req.parse(buffer[:readToIndex])
		if err != nil {
			fmt.Printf("Unexpedte error has occured while parsing data: %s", err)
			return &Request{}, err
		}

		newBuffer := make([]byte, len(buffer))

		copy(newBuffer, buffer[count:readToIndex])

		buffer = newBuffer

		readToIndex -= count

	}

	return &req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {

	returnRequest := Request{}

	consumedBytes := 0

	splitString := strings.Split(string(data), "\r\n")

	if len(splitString) == 1 {
		return &RequestLine{}, 0, nil
	}

	splitData := strings.Split(splitString[0], " ")

	if len(splitData) != 3 {
		fmt.Printf("Unexpected error has occured while spliting the data: %s", errors.New("Out of inedx, provided data is in a incorrect format"))
		return &RequestLine{}, -1, errors.New("Out of inedx, provided data is in a incorrect format")
	}

	versionParts := strings.Split(splitData[2], "/")

	if len(versionParts) != 2 || string(versionParts[0]) != "HTTP" {
		return &RequestLine{}, -1, errors.New("invalid HTTP version format")
	}

	if string(versionParts[1]) != "1.1" {
		return &RequestLine{}, -1, errors.New("unsupported HTTP version")
	}

	returnRequest.RequestLine.HttpVersion = string(versionParts[1])
	returnRequest.RequestLine.RequestTarget = string(splitData[1])
	returnRequest.RequestLine.Method = string(splitData[0])

	consumedBytes += len(splitString[0]) + 2

	for _, c := range returnRequest.RequestLine.Method {
		if !unicode.IsLetter(c) || !unicode.IsUpper(c) {
			return &RequestLine{}, -1, errors.New("invalid method: must be uppercase alphabetic characters")
		}
	}

	return &returnRequest.RequestLine, consumedBytes, nil

}

func (r *Request) parse(data []byte) (int, error) {

	counter := 0

	if r.state == 0 {

		res, count, err := parseRequestLine(data)

		if err != nil {
			return counter, err
		}

		if count == 0 && err == nil {
			return 0, nil
		}

		r.RequestLine = *res

		r.state = 1

		counter += count

		return counter, nil
	}

	if r.state == 1 {
		return 0, errors.New("Error: Trying to read data in a done state: %s")
	}

	return 0, errors.New("Error: Unkonwn state")

}
