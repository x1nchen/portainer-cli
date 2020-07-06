package unittest

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesContains(t *testing.T) {
	k := []byte("trade-agent-mt4client#61\n")
	sub := []byte("mt4cli")
	assert.True(t, bytes.Contains(k, sub))

	k = []byte("zzztest1\n")
	sub = []byte("test")
	assert.True(t, bytes.Contains(k, sub))
}
