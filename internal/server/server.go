package server

import (
	"fmt"
	R "httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool // async-safe flag, default is false
}

func Serve(port int) (*Server, error) {
	lsnr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener: lsnr}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	s.closed.Store(true)
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Fatalf("Error accepting connection: %s", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := R.WriteStatusLine(conn, R.OK)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	h := R.GetDefaultHeaders(0)
	err = R.WriteHeaders(conn, h)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
