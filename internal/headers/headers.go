package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	index := bytes.Index(data, []byte("\r\n"))

	if index == -1 {
		return 0, false, nil
	}

	if index == 0 {
		return 2, true, nil
	}

	headerLine := data[:index]

	keyValuePair := strings.SplitN(string(headerLine), ":", 2)

	if keyValuePair[0] != strings.TrimSpace(keyValuePair[0]) {
		return 0, false, errors.New("Provided data is in an incorrect format")
	}

	for i, val := range keyValuePair {
		keyValuePair[i] = strings.TrimSpace(val)
	}

	valid := strings.IndexFunc(keyValuePair[0], func(r rune) bool {
		return !isValidChar(r)
	}) == -1

	if !valid {
		return 0, false, errors.New("Provided data contains an invalid character")
	}

	key := strings.ToLower(keyValuePair[0])
	value := keyValuePair[1]

	existing, ok := h[key]

	if ok {
		h[key] = existing + ", " + value
	} else {
		h[key] = value
	}

	return index + 2, false, nil

}

func (h Headers) Get(key string) (string, bool) {

	val, ok := h[strings.ToLower(key)]

	if !ok {
		return "", false
	}

	return val, true
}

func isValidChar(c rune) bool {
	return ((c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '!' ||
		c == '#' ||
		c == '$' ||
		c == '%' ||
		c == '&' ||
		c == '\'' ||
		c == '*' ||
		c == '+' ||
		c == '-' ||
		c == '.' ||
		c == '^' ||
		c == '_' ||
		c == '`' ||
		c == '|' ||
		c == '~')
}
