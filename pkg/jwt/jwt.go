package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	Username  string `json:"username"`
	TokenType string `json:"token-type"`
	jwt.StandardClaims
}
type Token struct {
	TokenString string // token
}

type Result struct {
	UserName string // 返回的username
	Status   string // 状态：
	// case 0: "" (success)
	// case 1: "Token has expired"
	// case 2: "Error parsing token:"
}

const TokenTime = time.Hour * 24

var Secret = []byte("Author:Eutop1a")

const Author = "Eutop1a"

func GenerateToken(username string) string {
	c := MyClaims{
		username, "AccessToken", jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenTime).Unix(),
			Issuer:    Author,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	res, _ := token.SignedString(Secret)
	return res
}

func ParseToken() {

}
