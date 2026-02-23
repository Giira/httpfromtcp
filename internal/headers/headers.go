package headers

import (
	"bytes"
	"fmt"
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
			return 0, false, fmt.Errorf("error: parts hould only have 2 fields: %v", parts)
		}

		field_name := string(parts[0])
		if field_name != strings.TrimRight(field_name, " ") {
			return 0, false, fmt.Errorf("error: invalid header format - space before colon: %v", field_name)
		}
		field_name = strings.TrimSpace(field_name)

		field_value := strings.TrimSpace(string(parts[1]))
		h[field_name] = field_value
		return idx + 2, false, nil
	}

}

func NewHeaders() Headers {
	return Headers{}
}
