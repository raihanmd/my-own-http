package request

import (
	"errors"
	"io"
	"strings"
)

const initialBufferSize = 8

type parserState int

const (
	stateInitialized parserState = iota
	stateDone
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, initialBufferSize)
	readToIndex := 0
	request := &Request{state: stateInitialized}

	for request.state != stateDone {
		// Grow buffer if full
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// Read from reader
		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			request.state = stateDone
			break
		}
		if err != nil {
			return nil, err
		}

		readToIndex += n

		// Parse the data we have
		consumed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Remove consumed data from buffer
		if consumed > 0 {
			copy(buf, buf[consumed:readToIndex])
			readToIndex -= consumed
		}
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		line, consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if consumed == 0 {
			return 0, nil // Need more data
		}
		r.RequestLine = *line
		r.state = stateDone
		return consumed, nil

	case stateDone:
		return 0, errors.New("error: trying to read data in a done state")

	default:
		return 0, errors.New("error: unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// Find the end of the request line
	endIndex := -1
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			endIndex = i + 2
			break
		}
	}

	// If we don't have a complete line yet
	if endIndex == -1 {
		return nil, 0, nil
	}

	// Parse the line
	line := string(data[:endIndex])
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, 0, errors.New("invalid request line format")
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.Trim(strings.TrimPrefix(parts[2], "HTTP/"), "\r\n"),
	}, endIndex, nil
}
