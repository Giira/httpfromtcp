package request

import (
	"fmt"
	"io"
	"strings"
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
	rsInitialised requestState = iota
	rsDone
)

const bufferSize = 2

func RequestFromReader(reader io.Reader) (*Request, error) {
	b := make([]byte, bufferSize)
	readTo := 0

	r := &Request{
		state: rsInitialised,
	}
	for r.state != rsDone {
		if readTo >= len(b) {
			bNew := make([]byte, (len(b) * 2))
			copy(bNew, b)
			b = bNew
		}
		bytesRead, err := reader.Read(b[readTo:])
		if err != nil {
			if err == io.EOF {
				r.state = rsDone
				break
			}
			return nil, fmt.Errorf("error: failure to read from reader: %v", err)
		}

		readTo += bytesRead

		parsedTo, err := r.parse(b[:readTo])
		if err != nil {
			return nil, fmt.Errorf("error: failure to parse")
		}

		copy(b, b[parsedTo:])
		readTo -= parsedTo
	}
	return r, nil
}

func parseRequestLine(input []byte) (RequestLine, int, error) {
	bytes := 0
	lines := strings.Split(string(input), "\r\n")
	if len(lines) == 1 {
		return RequestLine{}, 0, nil
	}
	line := strings.Split(lines[0], " ")
	if len(line) != 3 {
		return RequestLine{}, bytes, fmt.Errorf("error: request line should always have 3 parts, not %v", len(line))
	}
	version := strings.Split(line[2], "/")
	rl := RequestLine{
		HttpVersion:   version[1],
		RequestTarget: line[1],
		Method:        line[0],
	}
	if version[0] != "HTTP" || version[1] != "1.1" {
		return RequestLine{}, bytes, fmt.Errorf("error: unrecognised http version")
	}
	if strings.ToUpper(rl.Method) != rl.Method {
		return RequestLine{}, bytes, fmt.Errorf("error: method not correctly formatted")
	}
	for _, char := range rl.Method {
		if char < 'A' || char > 'Z' {
			return RequestLine{}, bytes, fmt.Errorf("error: method contains non alphanumeric characters")
		}
	}
	return rl, bytes, nil
}

// Returns number of bytes parsed
func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case rsInitialised:
		rl, i, err := parseRequestLine(data)
		if err != nil {
			return i, err
		}
		if i == 0 {
			return 0, nil
		} else {
			r.RequestLine = rl
			r.state = rsDone
			return i, nil
		}
	case rsDone:
		return 0, fmt.Errorf("error: trying to read data in state: Done")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
