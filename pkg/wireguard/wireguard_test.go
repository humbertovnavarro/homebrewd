package wireguard

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKeypair(t *testing.T) {
	pub, priv := GenerateKeypair()
	assert.Equal(t, true, strings.HasSuffix(pub, "="), pub)
	assert.Equal(t, 44, len(pub))

	assert.Equal(t, true, strings.HasSuffix(priv, "="), priv)
	assert.Equal(t, 44, len(priv))
}
