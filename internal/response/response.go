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

type WriterState int

const (
	WriterInitialised WriterState = iota
	WriterStatusLineDone
	WriterHeadersDone
	WriterBodyDone
)

type Writer struct {
	state WriterState
	conn  io.Writer
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		conn:  conn,
		state: WriterInitialised,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriterInitialised {
		return fmt.Errorf("error: invalid state for writing statusline: %v", w.state)
	}

	var reasonString string

	switch statusCode {
	case 200:
		reasonString = "OK"
	case 400:
		reasonString = "Bad Request"
	case 500:
		reasonString = "Internal Server Error"
	default:
		reasonString = ""
	}

	_, err := io.WriteString(w.conn, fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonString))

	if err != nil {
		return fmt.Errorf("error: failed to write statusline: %v", err)
	}

	w.state = WriterStatusLineDone

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	out := headers.NewHeaders()
	out["Content-Length"] = fmt.Sprintf("%d", contentLen)
	out["Connection"] = "close"
	out["Content-Type"] = "text/plain"
	return out
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WriterStatusLineDone {
		return fmt.Errorf("error: invalid state for writing headers: %v", w.state)
	}

	out := []byte{}
	for k, v := range headers {
		out = fmt.Appendf(out, "%v: %v\r\n", k, v)
	}
	out = append(out, []byte("\r\n")...)
	_, err := w.conn.Write(out)
	if err != nil {
		return fmt.Errorf("error: failed to write headers to connection: %v", err)
	}

	w.state = WriterHeadersDone

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WriterHeadersDone {
		return 0, fmt.Errorf("error: incorrect writer state for body writing: %v", w.state)
	}
}
