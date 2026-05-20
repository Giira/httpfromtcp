package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget

	if strings.HasPrefix(target, "/httpbin") {
		handlerProxyHttpbin(w, req)
		return
	}

	var body []byte
	switch target {
	case "/yourproblem":
		w.WriteStatusLine(response.CodeBadRequest)
		body = []byte("<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>")
	case "/myproblem":
		w.WriteStatusLine(response.CodeServerError)
		body = []byte("<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>")
	default:
		w.WriteStatusLine(response.CodeOK)
		body = []byte("<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>")
	}
	h := response.GetDefaultHeaders(len(body))
	h.Change("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handlerProxyHttpbin(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	target = strings.TrimPrefix(target, "/httpbin/")
	url := "https://httpbin.org/" + target
	fmt.Printf("Attempting to send data to %v", url)

	res, err := http.Get(url)
	if err != nil {
		w.WriteStatusLine(response.CodeServerError)
		body := []byte("<html><head><title>500 Server Error</title></head><body><h1>Server Error</h1><p>Request to httpbin.org failed</p></body></html>")
		h := response.GetDefaultHeaders(len(body))
		h.Change("Content-Type", "`text/html`")
		w.WriteHeaders(h)
		w.WriteBody(body)
		return
	}

	defer res.Body.Close()

	w.WriteStatusLine(response.CodeOK)
	h := headers.NewHeaders()
	h.Set("Content-Type", res.Header.Get("Content-Type"))
	h.Set("Connection", "close")

	w.WriteHeaders(h)

	buffer := make([]byte, 1024)
	var unchunkedBody []byte

	for {
		n, err := res.Body.Read(buffer)
		unchunkedBody = append(unchunkedBody, buffer[:n]...)

		log.Printf("%v bytes read\n", n)

		if n > 0 {
			_, err := w.WriteChunkedBody(buffer[:n])
			if err != nil {
				log.Printf("error writing chunkedly: %v\n", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error: failure to read response body: %v", err)
			break
		}
	}

	bodyHash := sha256.Sum256(unchunkedBody)

	trailers := headers.NewHeaders()
	trailers.Set("X-Content-Sha256", string(fmt.Sprintf("%x", bodyHash)))
	trailers.Set("X-Content-Length", strconv.Itoa(len(unchunkedBody)))

	w.WriteTrailers(trailers)
}
