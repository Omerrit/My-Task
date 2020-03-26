package cookies

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func sign(sessionId []byte, secretKey string, duration time.Duration) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(duration).Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        string(sessionId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

func parse(token string, secretKey string) ([]byte, int64, error) {
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, 0, ErrIncorrectSessionToken
	}
	if !tkn.Valid {
		return nil, 0, ErrIncorrectSessionToken
	}
	if claims.ExpiresAt-time.Now().Unix() <= 0 {
		return nil, 0, ErrTokenExpired
	}

	return []byte(claims.Id), claims.ExpiresAt, nil
}
