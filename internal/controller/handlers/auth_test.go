package handlers

import (
	"testing"
	"time"

	"github.com/bbt-t/lets-go-keep/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticatorJWT(t *testing.T) {
	auth := newAuthenticatorJWT([]byte("secret_key"), time.Now().Add(1*time.Hour).Unix())
	assert.NotEmpty(t, auth)
}

func TestAuthenticatorJWT(t *testing.T) {
	auth := newAuthenticatorJWT([]byte("secret_key"), time.Now().Add(1*time.Hour).Unix())

	userID := entity.UserID("user_id_12")

	token, err := auth.CreateToken(userID)
	assert.NoError(t, err)

	id, errValidate := auth.ValidateToken(token)
	assert.NoError(t, errValidate)
	assert.Equal(t, userID, id)
}
