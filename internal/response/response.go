package response

import (
	H "httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

const CRLF = "\r\n"

func WriteStatusLine(w io.Writer, statuscode StatusCode) error {
	reasonPhrase := ""
	switch statuscode {
	case OK:
		reasonPhrase = "OK"
	case BadRequest:
		reasonPhrase = "Bad Request"
	case InternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	_, err := w.Write([]byte("HTTP/1.1 " + strconv.Itoa(int(statuscode)) + " " + reasonPhrase + CRLF))
	return err
}

func GetDefaultHeaders(contentLen int) H.Headers {
	h := H.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers H.Headers) error {
	for k, _ := range headers {
		v, _ := headers.Get(k)
		_, err := w.Write([]byte(k + ": " + v + CRLF))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte(CRLF))
	return err
}
