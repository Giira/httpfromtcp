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

const bufferSize = 2

func RequestFromReader(reader io.Reader) (*Request, error) {
	b := make([]byte, bufferSize)
	readTo := 0
	var r *Request
	r = &Request{
		State: 0,
	}
	for {
		if r.State == 0 {
			if readTo == len(b) {
				bNew := make([]byte, (readTo * 2))
				copy(bNew, b)
				b = bNew
			}
			bytesRead, err := reader.Read(b[readTo:])
			if err != nil {
				if err == io.EOF {
					if bytesRead != 0 {
						return nil, fmt.Errorf("error: if err is EOF, no bytes should have been read: %v", err)
					}
					if r.State != 1 {
						return nil, fmt.Errorf("error: Incomplete request: %v", err)
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

		}
	}
	rl, b, err := parseRequestLine(input)
	if err != nil {
		return nil, fmt.Errorf("error parsing request line: %v", err)
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
		if err != nil {
			return i, err
		}
		if i == 0 {
			return 0, nil
		} else {
			r.RequestLine = rl
			r.State = 1
			return i, nil
		}
	} else if r.State == 1 {
		return 0, fmt.Errorf("error: trying to read data in state: Done")
	} else {
		return 0, fmt.Errorf("error: unknown state")
	}
}
