package server

import (
	"bytes"
	"fmt"
	Req "httpfromtcp/internal/request"
	Resp "httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

const CRLF = "\r\n"

type Handler func(w io.Writer, req *Req.Request) *HandlerError

type HandlerError struct {
	StatusCode Resp.StatusCode
	Message    string
}

// copied from assignment, better DRY
func WriteHandlerErrorToBuffer(w io.Writer, he *HandlerError) {
	Resp.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := Resp.GetDefaultHeaders(len(messageBytes))
	Resp.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

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
		WriteHandlerErrorToBuffer(conn, hErr)
		return
	}

	Resp.WriteStatusLine(conn, Resp.OK)
	respHeaders := Resp.GetDefaultHeaders(buff.Len())
	err = Resp.WriteHeaders(conn, respHeaders)
	conn.Write(buff.Bytes())
}
