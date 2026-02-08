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

func PrintRequestLine(req *Request) {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %v\n", req.RequestLine.Method)
	fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
	fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)
}

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

func parseRequestLine(input []byte) (*RequestLine, int, error) {
	lines := strings.Split(string(input), "\r\n")
	if len(lines) == 1 {
		return nil, 0, nil
	}
	text := string(lines[0])
	rl, err := parseString(text)
	if err != nil {
		return nil, 0, err
	}
	return rl, len(text) + 2, nil
}

func parseString(str string) (*RequestLine, error) {
	sections := strings.Split(str, " ")
	if len(sections) != 3 {
		return nil, fmt.Errorf("error: request line should always have 3 parts, not %v - %v", len(sections), sections)
	}

	version := strings.Split(sections[2], "/")
	if version[0] != "HTTP" || version[1] != "1.1" {
		return nil, fmt.Errorf("error: unrecognised http version: %v", sections[2])
	}

	target := sections[1]

	method := sections[0]
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("error: method should be upper case letters only: %v", method)
		}
	}

	rl := &RequestLine{
		HttpVersion:   version[1],
		RequestTarget: target,
		Method:        method,
	}

	return rl, nil
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
			r.RequestLine = *rl
			r.state = rsDone
			return i, nil
		}
	case rsDone:
		return 0, fmt.Errorf("error: trying to read data in state: Done")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
