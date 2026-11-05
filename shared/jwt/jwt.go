package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mohammad")

type CustomClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`

	jwt.RegisteredClaims
}

func generateToken(userID string, email string, expireDelay time.Duration) (string, error) {
	expirationTime := time.Now().Add(expireDelay)

	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hotel-backend-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("Unauthorized")
		}

		// Return the secret key used to verify the signature
		return secretKey, nil
	})

	if err != nil {
		// return nil, err
		return nil, errors.New("Unauthorized")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Unauthorized")
}
