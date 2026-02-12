package headers

import "strings"

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	parts := strings.Split(string(data), "\r\n")

}

func NewHeaders() Headers {
	return Headers{}
}
