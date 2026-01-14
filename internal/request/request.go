package request

import (
	"io"
	"log"
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
		log.Fatalf("error reading request: %v", err)
	}
	rl := parseRequestLine(input)
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
