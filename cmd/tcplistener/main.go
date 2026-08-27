package main

import (
	"fmt"
	Req "httpfromtcp/internal/request"
	"log"
	"net"
)

const PORT = ":42069"

func main() {
	lsnr, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Error listening to network: %s", err)
	}
	defer lsnr.Close()

	fmt.Println("Listening for TCP traffic on", PORT)
	for {
		conn, err := lsnr.Accept()
		if err != nil {
			log.Fatalf("Error awaiting connection: %s", err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())
		reqChan, err := Req.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Error reading request: %s", err)
		}
		req := <-reqChan
		if req != nil {
			fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
				req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
