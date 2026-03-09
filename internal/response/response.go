package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	codeOK          StatusCode = 200
	codeBadRequest  StatusCode = 400
	codeServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case 200:
		w.Write([]byte("HTTP/1.1 200 OK"))
	case 400:
		w.Write([]byte("HTTP/1.1 400 Bad Request"))
	case 500:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error"))
	default:
		out := fmt.Sprintf("HTTP/1.1 %v", statusCode)
		w.Write([]byte(out))
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	out := headers.NewHeaders()
	out["Content-Length"] = strconv.Itoa(contentLen)
	out["Connection"] = "close"
	out["Content-Type"] = "text/plain"
	return out
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for header := range headers {
		out := fmt.Sprintf("%v: %v", header, headers[header])
		_, err := w.Write([]byte(out))
		if err != nil {
			return fmt.Errorf("error: failed to write headers: %v", err)
		}
	}
	w.Write([]byte("\r\n\r\n"))
	return nil
}
