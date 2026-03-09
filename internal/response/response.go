package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int

const (
	CodeOK          StatusCode = 200
	CodeBadRequest  StatusCode = 400
	CodeServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case 200:
		_, err = w.Write([]byte("HTTP/1.1 200 OK"))
	case 400:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request"))
	case 500:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error"))
	default:
		out := fmt.Sprintf("HTTP/1.1 %v", statusCode)
		_, err = w.Write([]byte(out))
	}
	if err != nil {
		return fmt.Errorf("error: failed to write statusline: %v", err)
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	out := headers.NewHeaders()
	out["Content-Length"] = fmt.Sprintf("%d", contentLen)
	out["Connection"] = "close"
	out["Content-Type"] = "text/plain"
	return out
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	out := []byte{}
	for k, v := range headers {
		out = fmt.Appendf(out, "%v: %v\r\n", k, v)
	}
	out = append(out, []byte("\r\n")...)
	_, err := w.Write(out)
	return err
}
