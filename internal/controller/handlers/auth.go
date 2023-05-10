package handlers

import (
	"log"

	"github.com/bbt-t/lets-go-keep/internal/entity"
	"github.com/bbt-t/lets-go-keep/internal/storage"

	"github.com/golang-jwt/jwt/v5"
)

// authenticatorJWT is authenticator which uses JWT.
type authenticatorJWT struct {
	secretKey      []byte
	expirationTime int64
}

// NewAuthenticatorJWT gets new authenticatorJWT.
func newAuthenticatorJWT(secretKey []byte, expirationTime int64) *authenticatorJWT {
	return &authenticatorJWT{
		secretKey:      secretKey,
		expirationTime: expirationTime,
	}
}

// CreateToken implementation of Authenticator interface. Creates token, which stores userID.
func (a *authenticatorJWT) CreateToken(userID entity.UserID) (entity.AuthToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"], claims["userID"] = a.expirationTime, userID

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		log.Println("Failed generate token for authentication:", err)
		return "", storage.ErrUnknown
	}

	return entity.AuthToken(tokenString), nil
}

// ValidateToken implementation of Authenticator interface. Validates token, returns userID.
func (a *authenticatorJWT) ValidateToken(token entity.AuthToken) (entity.UserID, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(string(token), claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, storage.ErrUnknown
		}
		return a.secretKey, nil
	})

	if err != nil {
		return "", storage.ErrUserUnauthorized
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", storage.ErrUserUnauthorized
	}

	return entity.UserID(userID), nil
}
