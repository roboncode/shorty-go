package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounter_EncodeCounter(t *testing.T) {
	counter := Counter{1}
	encoded := counter.EncodeCounter()

	assert.Equal(t, string(encoded), `{"value":1}`, "EncodeCounter failed to match")
}

func TestDecodeCounter(t *testing.T) {
	counter := Counter{1}
	encoded := counter.EncodeCounter()
	decoded, _ := DecodeCounter(encoded)
	assert.Equal(t, decoded.Value, 1, "DecodeCounter failed to match")
}
