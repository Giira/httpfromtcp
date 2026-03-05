package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"strconv"
)

type StatusCode int

const (
	200 StatusCode = iota
	400
	500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case 200:
		w.Write("HTTP/1.1 200 OK")
	case 400:
		w.Write("HTTP/1.1 400 Bad Request")
	case 500:
		w.Write("HTTP/1.1 500 Internal Server Error")
	default:
		out := fmt.Sprintf("HTTP/1.1 %v", statusCode)
		w.Write(out)
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	out := headers.NewHeaders()
	out["Content-Length"] = strconv.Atoi(contentLen)
	out["Connection"] = "close"
	out["Content-Type"] = "text/plain"
	return out
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	err := w.Write(headers)
	if err != nil {
		return fmt.Errorf("error: failed to write headers: %v", err)
	}
}
