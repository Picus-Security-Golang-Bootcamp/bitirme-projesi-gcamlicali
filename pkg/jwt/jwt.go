package jwt_helper

import (
	"encoding/json"
	"log"

	"github.com/golang-jwt/jwt"
)

type DecodedToken struct {
	Iat     int    `json:"iat"`
	IsAdmin bool   `json:"roles"`
	UserId  int    `json:"userId"`
	Email   string `json:"email"`
	Iss     string `json:"iss"`
}

func GenerateToken(claims *jwt.Token, secret string) string {
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)
	token, _ := claims.SignedString(hmacSecret)

	return token
}

func VerifyToken(token string, secret string) *DecodedToken {
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)

	decoded, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		log.Println("Decoded Hatasi: ", err)
		return nil
	}

	if !decoded.Valid {
		log.Println("Valid Hatasi")
		return nil
	}

	decodedClaims := decoded.Claims.(jwt.MapClaims)
	var decodedToken DecodedToken
	jsonString, _ := json.Marshal(decodedClaims)
	json.Unmarshal(jsonString, &decodedToken)

	return &decodedToken
}
