package headers

import (
	"bytes"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	switch idx {
	case -1:
		return 0, false, nil
	case 0:
		return 2, true, nil
	default:
		parts := bytes.SplitN(data[:idx], []byte(":"), 2)
		h[string(parts[0])] = string(parts[1])
		return len(parts) + 2, true, nil
	}

}

func NewHeaders() Headers {
	return Headers{}
}
