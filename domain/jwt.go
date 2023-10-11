// julo/domain/jwt.go
package domain

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// Claims represents the JWT claims structure.
type Claims struct {
	CustomerXID string `json:"customer_xid"`
	jwt.StandardClaims
}

// NewJWTClaims initializes a new JWT claims struct with user ID.
func NewJWTClaims(customer_xid string) *Claims {
	return &Claims{
		CustomerXID: customer_xid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(), // Set token expiration time
		},
	}

}

// GenerateJWTToken generates a JWT token with the given customer ID
func GenerateJWTToken(customerID string) (string, error) {
	// Access the JWT secret key from the configuration
	jwtSecret := os.Getenv("JWT_SECRET")
	// Create a new JWT token with a custom claim (in this case, "customer_id")
	claims := NewJWTClaims(customerID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
