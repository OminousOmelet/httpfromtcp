package main

import (
	Req "httpfromtcp/internal/request"
	Srvr "httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func handler(w io.Writer, req *Req.Request) *Srvr.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &Srvr.HandlerError{StatusCode: 400, Message: "Your problem is not my problem\n"}
	case "/myproblem":
		return &Srvr.HandlerError{StatusCode: 500, Message: "Woopsie, my bad\n"}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
}

func main() {
	server, err := Srvr.Serve(handler, port)
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
