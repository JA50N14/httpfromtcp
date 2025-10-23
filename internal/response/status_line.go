package response

import (
	"fmt"
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("writer is in wrong state: %d", w.writerState)
	}

	defer func() {w.writerState = writerStateHeaders}()

	var reasonPhrase string
	switch statusCode {
	case StatusCodeSuccess:
		reasonPhrase = "OK"
	case StatusCodeBadRequest:
		reasonPhrase = "Bad Request"
	case StatusCodeInternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
		return fmt.Errorf("invalid status code: %d", statusCode)
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase) 
	_, err := w.writer.Write([]byte(statusLine))
	return err
}