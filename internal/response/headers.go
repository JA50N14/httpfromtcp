package response

import (
	"fmt"
	"strings"

	"github.com/JA50N14/httpfromtcp/internal/headers"
)

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("writer is in wrong state: %d", w.writerState)
	}

	defer func() { w.writerState = writerStateBody }()

	for k, v := range h {
		key := strings.ToLower(k)
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, v)))
		if err != nil {
			return fmt.Errorf("error writing response headers to connection: %v", err)
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
