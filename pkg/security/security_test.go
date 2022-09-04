package security_test

import (
	"authstore/pkg/security"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "fuuuck2222"
	cost := 14
	hash, err := security.HashPassword(password, cost)

	assert.Nil(t, err)
	assert.True(t, security.CheckPasswordHash(password, hash))
}
