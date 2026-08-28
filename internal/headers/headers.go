package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)
	if !strings.Contains(str, "\r\n") {
		return 0, false, nil
	}
	if strings.HasPrefix(str, "\r\n") {
		return 0, true, nil
	}
	if strings.HasPrefix(str, " ") {
		return 0, false, errors.New("field line contains leading spaces")
	}
	str, _, _ = strings.Cut(str, "\r\n")
	bytesUsed := len(str) + 2 // +2 for '\r\n'
	key, value, colonFound := strings.Cut(str, ":")
	if !colonFound {
		return 0, false, errors.New("invalid format")
	}
	if strings.HasSuffix(key, " ") {
		return 0, false, errors.New("invalid spacing")
	}
	h[key] = strings.Trim(value, " ")
	//bytesTotal := crflBytes + len(key) + len(h[key])
	return bytesUsed, false, nil
}
