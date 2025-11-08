package middleware

import (
	"context"
	"net/http"
	"strings"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// защита от несанкционированного доступа
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "invalid token", http.StatusUnauthorized)
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			userID := claims["sub"].(string)
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateJWT(tokenStr string, jwtSecret string) (string, error) {
    if tokenStr == "" {
        return "", errors.New("empty token")
    }

    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })

    if err != nil || !token.Valid {
        return "", errors.New("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", errors.New("invalid token claims")
    }

    userID, ok := claims["sub"].(string)
    if !ok || userID == "" {
        return "", errors.New("userID missing in token")
    }

    return userID, nil
}