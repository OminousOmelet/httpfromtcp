package server

import (
	"bytes"
	"fmt"
	Req "httpfromtcp/internal/request"
	Resp "httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

const CRLF = "\r\n"

type Server struct {
	listener net.Listener
	closed   atomic.Bool // async-safe flag, default is false
	handler  Handler
}

func Serve(handler Handler, port int) (*Server, error) {
	lsnr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener: lsnr, handler: handler}
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
	req, err := Req.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error parsing request: %s", err)
	}
	var buff bytes.Buffer

	// special handler error...handling
	hErr := s.handler(&buff, req)
	if hErr != nil {
		err = WriteHandlerErrorToBuffer(&buff, hErr)
		if err != nil {
			fmt.Printf("failed to pass handler error: %s\n", err)
			return
		}
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d\r\n\r\n%s", hErr.StatusCode, hErr.Message)))
		return
	}
	msg := "All good, frfr\n"
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d\r\n\r\n%s", Resp.OK, msg)))
	// respHeaders := Resp.GetDefaultHeaders(0)
	// err = Resp.WriteStatusLine(conn, Resp.OK)
	// if err != nil {
	// 	fmt.Printf("error writing status: %s\n", err)
	// }
	// err = Resp.WriteHeaders(conn, respHeaders)
	// if err != nil {
	// 	fmt.Printf("error writing headers: %s\n", err)
	// }
	// err = Resp.WriteBody(conn, buff.Bytes())
	// if err != nil {
	// 	fmt.Printf("error writing body: %s\n", err)
	// }
}
