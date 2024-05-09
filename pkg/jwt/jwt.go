package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	UserIsAdmin bool   `json:"user_is_admin"`
	jwt.StandardClaims
}

const TokenTime = time.Hour * 24

var jwtSecret = []byte("Author:Eutop1a")

const Author = "Eutop1a"

// GenerateToken 签发用户token
func GenerateToken(id int64, username string, userIsAdmin bool) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(TokenTime)
	claims := Claims{
		UserID:      id,
		Username:    username,
		UserIsAdmin: userIsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    Author,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 验证用户token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
