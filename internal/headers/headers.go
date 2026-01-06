package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var HEADERS_SEP = []byte("\r\n")
var FIELD_VALUE_SEP = ":"

var ERROR_MALFORMED_HEADERS = fmt.Errorf("malformed headers")
var ERROR_MALFORMED_KEY = fmt.Errorf("malformed key")
var ERROR_INVALID_HEADER_KEY = fmt.Errorf("invalid char in header key")
var ERROR_MISSING_COLON = fmt.Errorf("missing colon")

const (
	specialChars = "!#$%&'*+-.^_`|~"
)

func isValidHeaderChar(r rune) bool {
	switch {
	case r >= 'A' && r <= 'Z':
		return true
	case r >= 'a' && r <= 'z':
		return true
	case r >= '0' && r <= '9':
		return true
	default:
		return strings.ContainsRune(specialChars, r)
	}
}

func isValidHeaderName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for _, r := range name {
		if !isValidHeaderChar(r) {
			return false
		}
	}
	return true
}

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) string {
	return h[key]
}

func (h Headers) Set(key string, val string) {
	h[key] = val
}

func (h Headers) Replace(key string, val string) {
	h[key] = val
}

func (h Headers) Parse(b []byte) (int, bool, error) {
	sepIdx := bytes.Index(b, HEADERS_SEP)

	if sepIdx == -1 {
		return 0, false, nil
	}

	if sepIdx == 0 {
		return len(HEADERS_SEP), true, nil
	}

	headerParts := strings.Split(string(b[:sepIdx]), FIELD_VALUE_SEP)

	if headerParts[0] != strings.TrimRight(headerParts[0], " ") {
		return 0, false, ERROR_MALFORMED_KEY
	}

	if len(headerParts) == 1 {
		return 0, false, ERROR_MISSING_COLON
	}

	if len(headerParts) != 2 {
		headerParts[1] = strings.Join(headerParts[1:], ":")
	}

	key := strings.ToLower(strings.TrimSpace(headerParts[0]))
	val := strings.TrimSpace(headerParts[1])

	if !isValidHeaderName(key) {
		return 0, false, ERROR_INVALID_HEADER_KEY
	}

	if mapVal, ok := h[key]; ok {
		val = strings.Join([]string{mapVal, val}, ", ")
	}

	h[key] = val
	return sepIdx + len(HEADERS_SEP), false, nil
}
