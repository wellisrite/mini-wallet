package middleware

import (
	"context"
	"julo/domain"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Middleware function for JWT authorization with "Token" prefix
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the Authorization header
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check if the header starts with "Token "
		if !strings.HasPrefix(authHeader, "Token ") {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Extract the JWT token value
		tokenString := strings.TrimPrefix(authHeader, "Token ")
		// Parse and verify the JWT token
		claims := &domain.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Token is not valid", http.StatusUnauthorized)
			return
		}

		// Extract the user ID from the claims and store it in the request context
		customerXID := claims.CustomerXID
		ctx := context.WithValue(r.Context(), "customerXID", customerXID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
