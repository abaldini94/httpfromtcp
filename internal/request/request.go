package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"httpfromtcp/internal/headers"
)

var supportedMethods = map[string]bool{
	"GET":    true,
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"UPDATE": true,
	"HEAD":   true,
}

var SEPARATOR = []byte("\r\n")
var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERROR_WRONG_REQUEST_LINE = fmt.Errorf("invalid number of elemenets in request-line")
var ERROR_INVALID_HTTP_VERSION = fmt.Errorf("malformed HTTP version specification")
var ERROR_UNSUPPORTED_HTTP_METHOD = fmt.Errorf("unsupported HTTP method")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported HTTP version")
var ERROR_MALFORMED_CONTENT_LENGTH = fmt.Errorf("content-length is not an integer")
var ERROR_BODY_GREATER_THAN_CONTENT_LENGTH = fmt.Errorf("body length is greater than content-length")

type parserState string

const (
	StateParsingInit    parserState = "init"
	StateParsingHeaders parserState = "headers"
	StateParsingBody    parserState = "body"
	StateParsingDone    parserState = "done"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine *RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parserState
}

func newRequest() *Request {
	return &Request{state: StateParsingInit, Headers: headers.NewHeaders()}
}

func (r *Request) parse(b []byte) (int, error) {

	switch r.state {
	case StateParsingDone:
		return 0, nil
	case StateParsingInit:
		rl, numByteParsed, err := parseRequestLine(b)
		if err != nil {
			return 0, err
		}

		if numByteParsed == 0 {
			return 0, nil
		}

		r.RequestLine = rl
		r.state = StateParsingHeaders
		return numByteParsed, nil
	case StateParsingHeaders:
		numBytesParsed, done, err := r.Headers.Parse(b)
		if err != nil {
			return 0, err
		}
		if numBytesParsed == 0 {
			return 0, nil
		}
		if done {
			r.state = StateParsingBody
		}
		return numBytesParsed, nil
	case StateParsingBody:
		cl := r.Headers.Get("content-length")
		if cl == "" {
			r.state = StateParsingDone
			return 0, nil
		}
		clen, err := strconv.Atoi(cl)
		if err != nil {
			return 0, ERROR_MALFORMED_CONTENT_LENGTH
		}
		for _, el := range b {
			r.Body = append(r.Body, el)
		}
		if len(r.Body) > clen {
			return 0, ERROR_BODY_GREATER_THAN_CONTENT_LENGTH
		}
		if len(r.Body) == clen {
			r.state = StateParsingDone
			return len(b), nil
		}
		return len(b), nil
	default:
		return 0, nil

	}
}

func (r *Request) done() bool {
	return r.state == StateParsingDone
}

func parseRequestLine(raw []byte) (*RequestLine, int, error) {
	sepIdx := bytes.Index(raw, SEPARATOR)

	if sepIdx == -1 {
		return nil, 0, nil
	}

	reqLine := raw[:sepIdx]
	read := sepIdx + len(SEPARATOR)
	parts := bytes.Split(reqLine, []byte(" "))

	if len(parts) != 3 {
		return nil, read, errors.New("Invalid Request Line")
	}

	method := string(parts[0])
	if strings.ToUpper(method) != method {
		return nil, read, ERROR_WRONG_REQUEST_LINE
	}

	if ok := supportedMethods[method]; !ok {
		return nil, read, ERROR_UNSUPPORTED_HTTP_METHOD
	}

	version := bytes.Split(parts[2], []byte("/"))
	if len(version) != 2 || string(version[0]) != "HTTP" {
		return nil, read, ERROR_INVALID_HTTP_VERSION
	}

	if string(version[1]) != "1.1" {
		return nil, read, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: string(parts[1]),
		HttpVersion:   string(version[1]),
	}, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// TODO edge case: reqLine/headers/body > 1K
	buf := make([]byte, 1024)
	bufLen := 0

	numBytesRecv, err := reader.Read(buf[bufLen:])
	if err != nil {
		return nil, err
	}
	bufLen += numBytesRecv
	for !request.done() {

		numBytesParsed, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:bufLen])
		bufLen -= numBytesParsed

		if numBytesParsed == 0 && !request.done() {
			numBytesRecv, err := reader.Read(buf[bufLen:])
			if err != nil {
				return nil, err
			}
			bufLen += numBytesRecv
		}
	}

	return request, nil
}
