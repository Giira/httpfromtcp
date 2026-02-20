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
		if field_name[0] == ' ' {
			return 0, false, fmt.Errorf("error: incorrect header format - leading space: %v", field_name)
		}
		if strings.Contains(field_name, " :") {
			return 0, false, fmt.Errorf("error: incorrect header format - space before colon: %v", field_name)
		}
		field_value := strings.TrimSpace(string(parts[1]))
		h[field_name] = field_value
		return len(field_name) + len(field_value) + 2, true, nil
	}

}

func NewHeaders() Headers {
	return Headers{}
}
