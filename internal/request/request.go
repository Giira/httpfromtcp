package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading from io reader: %v", err)
	}
	rl, b, err := parseRequestLine(input)
	if err != nil {
		return nil, fmt.Errorf("error parsing request line: %v", err)
	}
	var r *Request
	r = &Request{
		RequestLine: rl,
		State:       0,
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

func (r *Request) parse(data []byte) (int, error) {
	if r.State == 0 {
		rl, i, err := parseRequestLine(data)
	}
}
