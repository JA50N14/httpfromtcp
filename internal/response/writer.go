package response

import (
	"io"

)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	writerState writerState
	writer io.Writer
}

type StatusCode int
const (
	StatusCodeSuccess StatusCode = 200
	StatusCodeBadRequest StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)


func NewWriter(w io.Writer) *Writer {
	return &Writer {
		writerState: writerStateStatusLine,
		writer: w,
	}
}

