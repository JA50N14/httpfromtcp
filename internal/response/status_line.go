package response

import (
	"fmt"
	"net"

	"github.com/JA50N14/httpfromtcp/internal/headers"
)

type Writer struct {
	StatusCode StatusCode
	Headers headers.Headers
	BodyLen int
	Body []byte
	Out net.Conn
	WriterState WriterState
}

type StatusCode int
const (
	StatusCodeSuccess StatusCode = 200
	StatusCodeBadRequest StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

type WriterState int
const (
	StateWriteStatusLine WriterState = iota
	StateWriteHeaders
	StateWriteBody
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState != StateWriteStatusLine {
		return fmt.Errorf("write to status-line first")
	}
	
	var reasonPhrase string
	switch statusCode {
	case StatusCodeSuccess:
		reasonPhrase = "OK"
	case StatusCodeBadRequest:
		reasonPhrase = "Bad Request"
	case StatusCodeInternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
		return fmt.Errorf("invalid statusCode: %d", statusCode)
	}
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
	w.Out.Write([]byte(statusLine))
	w.WriterState = StateWriteHeaders
	return nil
}