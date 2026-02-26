package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte

	state   requestState
	bodyLen int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	rsInitialised requestState = iota
	rsParsingHeaders
	rsParsingBody
	rsDone
)

const bufferSize = 2

func PrintRequestLine(req *Request) {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %v\n", req.RequestLine.Method)
	fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
	fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)
}

func PrintHeaders(req *Request) {
	fmt.Println("Headers:")
	for key, value := range req.Headers {
		fmt.Printf("- %v: %v\n", key, value)
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b := make([]byte, bufferSize)
	readTo := 0

	r := &Request{
		state:   rsInitialised,
		Headers: headers.NewHeaders(),
		Body:    make([]byte, 0),
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
				if r.state != rsDone {
					return nil, fmt.Errorf("error: incomplete request")
				}
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

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case rsInitialised:
		rl, i, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if i == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.state = rsParsingHeaders
		return i, nil
	case rsParsingHeaders:
		i, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = rsParsingBody
		}
		return i, nil
	case rsParsingBody:
		cLength, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.state = rsDone
			return len(data), nil
		}

		cLengthI, err := strconv.Atoi(cLength)
		if err != nil {
			return 0, err
		}
		r.Body = append(r.Body, data...)
		r.bodyLen += len(data)

		if r.bodyLen > cLengthI {
			return 0, fmt.Errorf("error: body exceeds content length")
		}
		if r.bodyLen == cLengthI {
			r.state = rsDone
		}
		return len(data), nil
	case rsDone:
		return 0, fmt.Errorf("error: trying to read data in state: Done")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

// Returns number of bytes parsed
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != rsDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}
