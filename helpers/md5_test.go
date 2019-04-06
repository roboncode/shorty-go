package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMD5(t *testing.T) {
	hash := MD5("Hello, world!")
	assert.Equal(t, hash, "6cd3556deb0da54bca060b4c39479839", "MD5 hash failed to match")
}
