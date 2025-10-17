package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"strconv"

	"github.com/JA50N14/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body []byte
	bodyLengthRead int
	state       requestState

}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		Headers: headers.NewHeaders(),
		Body: make([]byte, 0),
		state:   requestStateInitialized,
	}

	for req.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.state, numBytesRead)
				}
				break
			}
			return nil, err
		}

		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[numBytesParsed:readToIndex])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("invalid method in request-line: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("invalid start-line: %s", str)
	}
	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("invalid HTTP-version: %s", str)
	}
	httpVersion := versionParts[1]
	if httpVersion != "1.1" {
		return nil, fmt.Errorf("invalid HTTP version: %s", httpVersion)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something went wrong
			return 0, err
		}
		if n == 0 {
			//just need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			//something went wrong
			return 0, err
		}
		if done {
			//end of Headers
			r.state = requestStateParsingBody
		}
		return n, nil
	case requestStateParsingBody:
		contentLenStr, ok := r.Headers.Get("Content-Length") 
		if !ok {
			//assuming if no content-length header is present, there is no body
			r.state = requestStateDone
			return len(data), nil
		}
		contentLen, err := strconv.Atoi(contentLenStr)
		if err != nil {
			return 0, fmt.Errorf("invalid content-length header: %s", err)
		}
		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)
		if r.bodyLengthRead > contentLen {
			return 0, fmt.Errorf("request body larger than content-length specified in header")
		}
		if r.bodyLengthRead == contentLen {
			r.state = requestStateDone
		}
		return len(data), nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}