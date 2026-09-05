package request

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const BUFFER_SIZE = 8

// faking enums is weird in Go
type Status int

const (
	Initialized = iota
	ParsingHeaders
	ParsingBody
	Done
)

type Request struct {
	Headers     headers.Headers
	RequestLine RequestLine
	Status      Status
	Body        []byte
	BodyLength  int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, BUFFER_SIZE)
	readToIndex := 0
	req := &Request{Headers: headers.NewHeaders(), Status: Initialized, Body: []byte{}, BodyLength: 0}
	for req.Status != Done {
		if readToIndex >= len(buff) {
			temp := make([]byte, len(buff)*2)
			copy(temp, buff)
			buff = temp
		}
		n, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if req.Status != Done {
					return nil, errors.New("incomplete request")
				}
				break
			}
			return nil, err
		}
		readToIndex += n
		n, err = req.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[n:])
		readToIndex -= n
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.Status != Done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	err := writeToDataLog(data)
	if err != nil {
		log.Printf("Error writing to data log: %s", err)
	}
	switch r.Status {
	case Initialized:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return n, err
		}
		if n == 0 {
			return 0, nil // keep trying until request line is complete
		}
		r.RequestLine = *reqLine
		r.Status = ParsingHeaders
		return n, nil
	case ParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.Status = ParsingBody
		}
		return n, nil
	case ParsingBody: // Copied from solution files, pffffff
		contentLenStr, ok := r.Headers.Get("Content-Length")
		if !ok {
			// assume that if no content-length header is present, there is no body
			r.Status = Done
			return len(data), nil
		}
		contentLen, err := strconv.Atoi(contentLenStr)
		if err != nil {
			return 0, fmt.Errorf("malformed Content-Length: %s", err)
		}
		r.Body = append(r.Body, data...)
		r.BodyLength += len(data)
		if r.BodyLength > contentLen {
			return 0, fmt.Errorf("Content-Length too large")
		}
		if r.BodyLength == contentLen {
			r.Status = Done
		}
		return len(data), nil
	case Done:
		return 0, errors.New("error during read attempt: data already in 'Done' status")
	default:
		return 0, errors.New("unknown status")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	str, _, complete := strings.Cut(string(data), "\r\n")
	if !complete {
		return nil, 0, nil // expected behavior until request line is complete
	}
	reqLine := &RequestLine{}
	reqLine.HttpVersion = str
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, 0, errors.New("wrong number of parts in request line")
	}
	for _, r := range parts[0] {
		if !unicode.IsUpper(r) {
			return nil, 0, errors.New("invalid http method")
		}
	}
	if parts[2] != "HTTP/1.1" {
		return nil, 0, errors.New("invalid http version")
	}
	reqLine.Method = parts[0]
	reqLine.RequestTarget = parts[1]
	reqLine.HttpVersion = parts[2][5:]
	return reqLine, len(str) + 2, nil
}

func writeToDataLog(data []byte) error {
	f, err := os.OpenFile("dataLog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
