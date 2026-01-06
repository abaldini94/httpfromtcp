package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int

const (
	StatusOK          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusServerError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("content-length", fmt.Sprintf("%d", contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	return h
}

type Writer struct {
	Conn io.Writer
}

func (w *Writer) WriteStatusLine(status StatusCode) error {
	switch status {
	case StatusOK:
		statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", StatusOK, "OK")
		_, err := w.Conn.Write([]byte(statusLine))
		if err != nil {
			return err
		}
		return nil
	case StatusBadRequest:
		statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", StatusBadRequest, "Bad Request")
		_, err := w.Conn.Write([]byte(statusLine))
		if err != nil {
			return err
		}
		return nil
	case StatusServerError:
		statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", StatusServerError, "Internal Server Error")
		_, err := w.Conn.Write([]byte(statusLine))
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for key, value := range headers {
		line := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Conn.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	_, err := w.Conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(b []byte) (int, error) {
	n, err := w.Conn.Write(b)
	return n, err
}
