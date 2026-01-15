package request

import (
	"fmt"
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
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading from io reader: %v", err)
	}
	rl := parseRequestLine(input)
	var r *Request
	r = &Request{
		RequestLine: rl,
	}
	return r, nil
}

func parseRequestLine(input []byte) RequestLine {
	lines := strings.Split(string(input), "\r\n")
	rl := RequestLine{
		HttpVersion:   lines[0],
		RequestTarget: lines[1],
		Method:        lines[2],
	}
	return rl
}
