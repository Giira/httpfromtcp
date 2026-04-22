package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	target = strings.TrimPrefix(target, "/httpbin")
	if target == "" {
		target = "/"
	}
	url := "https//httpbin.org" + target

	res, err := http.Get(url)
	var body []byte
	if err != nil {
		w.WriteStatusLine(response.CodeBadRequest)
		body = []byte("<html><head><title>500 Server Error</title></head><body><h1>Server Error</h1><p>Request to httpbin.org failed</p></body></html>")
		h := response.GetDefaultHeaders(len(body))
		h.Change("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody(body)
		return
	}

}
