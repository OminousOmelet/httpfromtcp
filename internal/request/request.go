package request

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

const BUFFER_SIZE = 8

// faking enums is weird in Go
type Status int

const (
	Initialized = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	Status      Status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (<-chan *Request, error) {
	reqChan := make(chan *Request)
	go func() {
		defer close(reqChan)
		buff := make([]byte, BUFFER_SIZE)
		readToIndex := 0
		req := &Request{Status: Initialized}
		for req.Status != Done {
			if readToIndex >= len(buff) {
				temp := make([]byte, len(buff)*2)
				copy(temp, buff)
				buff = temp
			}
			n, err := reader.Read(buff[readToIndex:])
			if err == io.EOF {
				req.Status = Done
				break // redundant but whatever
			}
			if err != nil {
				return
			}
			readToIndex += n
			n, err = req.parse(buff[:readToIndex])
			if err != nil {
				return
			}
			if n != 0 {
				temp := make([]byte, len(buff)-n)
				copy(temp, buff[n:])
				buff = temp
				readToIndex -= n
			}
		}
		reqChan <- req
	}()
	return reqChan, nil
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

func (r *Request) parse(data []byte) (int, error) {
	if r.Status == Done {
		return 0, errors.New("request already parsed")
	}
	if r.Status != Initialized {
		return 0, nil // ignore until initialized
	}
	err := writeToDataLog(data)
	if err != nil {
		log.Printf("Error writing to data log: %s", err)
	}
	reqLine, n, err := parseRequestLine(data)
	if err != nil {
		return n, err
	}
	if n == 0 {
		return 0, nil // keep trying until request line is complete
	}
	r.RequestLine = *reqLine
	r.Status = Done
	return n, nil
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
