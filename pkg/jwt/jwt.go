package jwt

import (
	"time"

	go_jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	go_jwt.StandardClaims
}

func SetJwtSecret(jwt string) {
	jwtSecret = []byte(jwt)
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		go_jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ginS",
		},
	}

	tokenClaims := go_jwt.NewWithClaims(go_jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := go_jwt.ParseWithClaims(token, &Claims{}, func(token *go_jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
