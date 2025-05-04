package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, _ := io.ReadAll(reader)

	result, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *result}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	requestStr := string(data)
	lines := strings.Split(requestStr, "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("invalid HTTP req")
	}

	parts := strings.Split(lines[0], " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid HTTP method")
	}

	return &RequestLine{Method: parts[0], RequestTarget: parts[1], HttpVersion: strings.TrimPrefix(parts[2], "HTTP/")}, nil
}
