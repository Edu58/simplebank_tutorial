package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(8)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	// Correct Password
	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	// Wrong Password
	err = CheckPassword("secret", hashedPassword)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword)
}
