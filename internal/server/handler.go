package server

import (
	"fmt"
	Req "httpfromtcp/internal/request"
	"io"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *Req.Request) *HandlerError

func WriteHandlerErrorToBuffer(w io.Writer, h *HandlerError) error {
	_, err := w.Write([]byte(fmt.Sprintf("%d %s", h.StatusCode, h.Message)))
	if err != nil {
		return fmt.Errorf("error: %v\n", err)
	}
	return nil
}
