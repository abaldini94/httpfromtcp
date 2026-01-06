package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	val, ok := headers["host"]
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069", val)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra withspaces
	headers = NewHeaders()
	data = []byte("   Host: localhost:42069    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 30, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("   Host: localhost:42069    \r\n\r\n")
	n, _, _ = headers.Parse(data)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, n, 2)

	// Test: Valid duplicated keys
	headers = NewHeaders()
	headers["set-person"] = "person1"
	data = []byte("Set-Person: person2\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "person1, person2", headers["set-person"])
	assert.Equal(t, 21, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid char in field name
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
