package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyLuhn(t *testing.T) {
	assert.True(t, VerifyLuhn("12345678903"))
	assert.False(t, VerifyLuhn("12345678904"))
}
