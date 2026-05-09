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
	WriterBodyInitialised
	WriterBodyDone
)

type Writer struct {
	state WriterState
	Conn  io.Writer
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		Conn:  conn,
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

	_, err := io.WriteString(w.Conn, fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonString))

	if err != nil {
		return fmt.Errorf("error: failed to write statusline: %v", err)
	}

	w.state = WriterStatusLineDone

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	out := headers.NewHeaders()
	out.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	out.Set("Connection", "close")
	out.Set("Content-Type", "text/plain")
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
	_, err := w.Conn.Write(out)
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

	n, err := w.Conn.Write(p)
	if err != nil {
		return 0, fmt.Errorf("error: failed to write body to connection: %v", err)
	}

	w.state = WriterBodyDone
	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state == WriterHeadersDone {
		w.state = WriterBodyInitialised
	}
	if w.state != WriterBodyInitialised {
		return 0, fmt.Errorf("error: incorrect state for chunked body writing: %v", w.state)
	}

	i := len(p)
	_, err := io.WriteString(w.Conn, fmt.Sprintf("%d\r\n", i))
	if err != nil {
		return 0, fmt.Errorf("error: failed to write %d to connection: %v", i, err)
	}

	n, err := w.Conn.Write(p)
	if err != nil {
		return 0, fmt.Errorf("error: failed to write: %v", err)
	}

	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != WriterBodyInitialised {
		return 0, fmt.Errorf("error: wrong writer state: %v", w.state)
	}
	w.state = WriterBodyDone

	n, err := io.WriteString(w.Conn, "0\r\n\r\n")
	if err != nil {
		return n, fmt.Errorf("error: failed to write: %v", err)
	}
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {

}
