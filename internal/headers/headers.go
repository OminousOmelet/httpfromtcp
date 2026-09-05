package headers

import (
	"errors"
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)
	if !strings.Contains(str, CRLF) {
		return 0, false, nil // possibly incomplete data?
	}
	if strings.HasPrefix(str, CRLF) {
		return 2, true, nil
	}
	if strings.HasPrefix(str, " ") {
		return 0, false, errors.New("field line contains leading spaces")
	}
	str, _, _ = strings.Cut(str, CRLF)
	bytesUsed := len(str) + 2 // +2 for CRLF ('\r\n')
	key, value, colonFound := strings.Cut(str, ":")
	if !colonFound {
		return 0, false, errors.New("invalid format")
	}
	if strings.HasSuffix(key, " ") {
		return 0, false, errors.New("invalid spacing")
	}
	if !isValidFieldName(key) {
		return 0, false, errors.New("invalid field name")
	}
	h.Set(key, value)
	return bytesUsed, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	value = strings.TrimSpace(value)
	// check if field exists, and add entry if not already present
	if _, ok := h[key]; ok {
		if !strings.Contains(h[key], value) {
			h[key] += ", " + value
		}
	} else {
		h[key] = value
	}
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	v, ok := h[key]
	return v, ok
}

var special_chars = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '^', '_', '`', '|', '~'}

func isValidFieldName(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			if !slices.Contains(special_chars, r) {
				return false
			}
		}
	}
	return true
}
