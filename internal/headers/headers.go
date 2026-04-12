package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
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
		if len(parts) != 2 {
			return 0, false, fmt.Errorf("error: parts should only have 2 fields: %v", parts)
		}

		field_name := string(parts[0])
		if field_name != strings.TrimRight(field_name, " ") {
			return 0, false, fmt.Errorf("error: invalid header format - space before colon: %v", field_name)
		}
		field_name = strings.TrimSpace(field_name)
		ok := checkChars(field_name)
		if !ok {
			return 0, false, fmt.Errorf("error: invalid character in header field_name: %v", field_name)
		}
		field_value := strings.TrimSpace(string(parts[1]))

		h.Set(field_name, field_value)
		return idx + 2, false, nil
	}

}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	out, ok := h[key]
	return out, ok
}

func (h Headers) Set(key string, value string) {
	key = strings.ToLower(key)
	current, ok := h[key]
	if ok {
		value = current + ", " + value
	}
	h[key] = value
}

func (h Headers) Change(key string, value string) {
	key = strings.ToLower(key)
	delete(h, key)
	h[key] = value
}

func checkChars(text string) bool {
	spec := strings.Split("!#$%&'*+-.^_`|~", "")
	for _, char := range text {
		if ('A' <= char && char <= 'Z') ||
			('a' <= char && char <= 'z') ||
			('0' <= char && char <= '9') ||
			slices.Contains(spec, string(char)) {
			continue
		} else {
			return false
		}
	}
	return true
}

func NewHeaders() Headers {
	return map[string]string{}
}
