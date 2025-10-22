package response

import (
	"fmt"

	"github.com/JA50N14/httpfromtcp/internal/headers"
)

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.WriterState != StateWriteHeaders {
		return fmt.Errorf("write to status-line before writing to headers")
	}
	
	h = AddDefaultHeaders(h, w.BodyLen)
	var headersString string
	for key, value := range h {
		headersString += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	headersString += "\r\n"
	_, err := w.Out.Write([]byte(headersString))
	if err != nil {
		return err
	}
	w.WriterState = StateWriteBody
	w.Headers = h
	return nil	
}


func AddDefaultHeaders(h headers.Headers, contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()
	defaultHeaders.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	defaultHeaders.Set("Connection", "close")
	defaultHeaders.Set("Content-Type", "text/plain")
	
	for key, value := range defaultHeaders {
		if _, ok := h[key]; ok {
			continue
		}
		h.Set(key, value)
	}
	return h
}
