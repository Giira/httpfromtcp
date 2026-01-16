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
	rl, err := parseRequestLine(input)
	if err != nil {
		return nil, fmt.Errorf("error parsing request line: %v", err)
	}
	var r *Request
	r = &Request{
		RequestLine: rl,
	}
	return r, nil
}

func parseRequestLine(input []byte) (RequestLine, error) {
	lines := strings.Split(string(input), "\r\n")
	line := strings.Split(lines[0], " ")
	if len(line) != 3 {
		return RequestLine{}, fmt.Errorf("error: request line should always have 3 parts, not %v", len(line))
	}
	rl := RequestLine{
		HttpVersion:   line[2],
		RequestTarget: line[1],
		Method:        strings.TrimLeft(line[0], "HTTP/"),
	}
	fmt.Println(rl.HttpVersion)
	return rl, nil
}
