package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("my-secret")

func GenerateToken(userID int32) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		},
	)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenstring string) (int32, error) {
	token, err := jwt.Parse(
		tokenstring,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	id := int32(claims["user_id"].(float64))
	return id, nil
}