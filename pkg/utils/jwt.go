package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("Something really secret")

type Claims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(id int64, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		Id:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Blue",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}
