package response

import (
	"fmt"
	"io"

	"github.com/JA50N14/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var headersString string
	for key, value := range headers {
		headersString += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	headersString += "\r\n"
	_, err := w.Write([]byte(headersString))
	return err
}
