package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndVerify(t *testing.T) {
	password := "my_secret_password"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, hash, password)

	assert.True(t, VerifyPassword(password, hash))
}
