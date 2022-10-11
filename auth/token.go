package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateAllTokens(email string) (signedToken string, err error) {
	expirationTime := time.Now().Local().Add(time.Hour * time.Duration(1))
	claims := &SignedDetails{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	return token, err
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}
