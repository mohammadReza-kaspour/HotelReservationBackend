package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mohammad")

type CustomClaims struct {
	UserID string `json:"user_id"`
	Role string `json:"role"`
	// RegisteredClaims embeds standard RFC 7519 fields (exp, iat, nbf, iss, etc.)
	jwt.RegisteredClaims
}

// GenerateToken creates and signs a new JWT token for a given user.
func GenerateToken(userID string, role string) (string, error) {
	// Define the expiration time (e.g., 2 hours from now)
	expirationTime := time.Now().Add(2 * time.Hour)

	// Populate custom claims together with standard registered claims
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// ExpiresAt dictates when the token becomes invalid
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// IssuedAt tracks the exact timestamp when this token was minted
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Issuer identifies the service that created this token
			Issuer: "hotel-backend-api",
			// Subject typically identifies the principal/user ID
			Subject: userID,
		},
	}

	// Create a new token object with the HMAC-SHA256 signing method and custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using our secret key and return the signed compact string (header.payload.signature)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		// Return empty string and error if signing fails
		return "", err
	}

	// Return the completed signed token string
	return tokenString, nil
}

// ValidateToken parses, validates the signature, checks expiration, and returns the claims.
func ValidateToken(tokenString string) (*CustomClaims, error) {
	// ParseWithClaims decodes the token string and populates our CustomClaims struct
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		// Verify that the token's signing algorithm matches what we expect (HMAC in this case)
		// This protects against "alg: none" or algorithm-switching vulnerabilities
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Reject tokens signed with unexpected methods (e.g., RSA when HMAC is expected)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key used to verify the signature
		return secretKey, nil
	})

	// Check if parsing encountered errors (such as malformed token or signature mismatch)
	if err != nil {
		return nil, err
	}

	// Extract the typed claims and verify that the token is marked valid by the parser
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Token is verified, not expired, and contains valid claims
		return claims, nil
	}

	// Return a generic invalid token error if validation checks failed
	return nil, errors.New("invalid or expired token")
}

func main() {
	// 1. Generate a token for a user
	token, err := GenerateToken("user_987654", "admin")
	if err != nil {
		// Terminate if generation encounters an issue
		panic(fmt.Sprintf("Failed to generate token: %v", err))
	}

	// Print the signed JWT compact string
	fmt.Printf("Generated JWT:\n%s\n\n", token)


	time.Sleep(time.Second * 5)
	// 2. Validate and decode the token
	claims, err := ValidateToken(token)
	if err != nil {
		// Handle validation failure (e.g., expired, tampered, or invalid signature)
		fmt.Printf("Token validation failed: %v\n", err)
		return
	}

	// Print extracted data to confirm successful decoding
	fmt.Println("Token successfully verified:")
	fmt.Printf("- User ID: %s\n", claims.UserID)
	fmt.Printf("- Role: %s\n", claims.Role)
	fmt.Printf("- Issuer: %s\n", claims.Issuer)
	fmt.Printf("- Expires At: %s\n", claims.ExpiresAt.Time)
}