package token

import (
	"github.com/stretchr/testify/assert"
	"reunion/configuration"
	"testing"
)

func TestCreateToken(test *testing.T) {
	configuration.Conf = new(configuration.TestConfiguration)
	token := Create("")

	assert.NotEmpty(test, token)
}

func TestVerifyToken(test *testing.T) {
	token := Create("")

	valid := Verify(token)

	assert.True(test, valid)
}

func TestVerifyMandatoryClaimsToken(test *testing.T) {
	token := Create("user_authentication")

	valid := Verify(token)

	assert.True(test, valid)
}

func TestVerifyInvalidToken(test *testing.T) {
	expiredToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MTM4MDQwMjZ9.tQc1j0lTED2yH2QYEtV6zzzsZ2bVVTNTpgF6gDl-lVEQEM2hOhJxpBDs-ulOXPUk3SoVm0UyxJAX9hspkb2A_Z7ioW8IKSnhWj9303Dhn9LY7ZAo3wMyyHKV4dR2xqqXcQJaSh9Kbk03mau0m3Kqtw-N1W45mEcxXXk3IzFbSMA"

	invalid := Verify(expiredToken)

	assert.False(test, invalid)
}
