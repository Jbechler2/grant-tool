package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyRole   contextKey = "role"
)

type JWTMiddleware struct {
	secret string
}

func NewJWTMiddleware(secret string) func(http.Handler) http.Handler {
	m := &JWTMiddleware{secret: secret}
	return m.verify
}

func (m *JWTMiddleware) verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)

		if tokenString == "" {
			writeUnauthorized(w)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.secret), nil
		})
		if err != nil || !token.Valid {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
				writeTokenExpired(w)
			} else {
				writeUnauthorized(w)
			}
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeUnauthorized(w)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			writeUnauthorized(w)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			writeUnauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
		ctx = context.WithValue(ctx, ContextKeyRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(string)
	return id, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyRole).(string)
	return role, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error": "unauthorized"}`))
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	cookie, err := r.Cookie("token")
	if err == nil {
		return cookie.Value
	}

	return ""

}

func writeTokenExpired(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error": "token_expired"}`))
}
