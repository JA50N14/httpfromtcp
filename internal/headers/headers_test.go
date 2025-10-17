package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestHeadersParse(t *testing.T) {
	//Test: Valid single header with done
	headers := NewHeaders()
	data := []byte("HoSt: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, n, 2)
	assert.Equal(t, true, done)

	//Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: Single header not done
	headers = NewHeaders()
	data = []byte("HosT: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.Equal(t, false, done)

	//Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("  Host:   localhost: 42069      \r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost: 42069", headers["host"])
	assert.Equal(t, 34, n)
	assert.Equal(t, false, done)

	//Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["content-type"] = "application/json"
	headers["content-length"] = "50"
	data = []byte("HOst: localhost:42069\r\nUser-Agent: curl/7.8.1\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 3, len(headers))
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.Equal(t, false, done)

	n, done, err = headers.Parse(data[23:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 4, len(headers))
	assert.Equal(t, "curl/7.8.1", headers["user-agent"])
	assert.Equal(t, 24, n)
	assert.Equal(t, false, done)

	//Test: Valid done
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.Equal(t, false, done)

	n, done, err = headers.Parse(data[23:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.Equal(t, 1, len(headers))
	assert.Equal(t, true, done)

	//Test: Invalid spacing headers
	headers = NewHeaders()
	data = []byte("Host : localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, len(headers))
	assert.Equal(t, 0, n)
	assert.Equal(t, false, done)

	//Test: Invalid character in header key
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, len(headers))
	assert.Equal(t, 0, n)
	assert.Equal(t, false, done)

	//Test: Valid starting header that matches the header in data to be parsed
	headers = NewHeaders()
	headers["host"] = "localhost:42069"
	data = []byte("Host: localhost:42070\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, 23, n)
	require.Equal(t, 1, len(headers))
	require.Equal(t, false, done)
}	
