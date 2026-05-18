package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ContextKeyTenantID contextKey = "tenant_id"
	ContextKeyActorID  contextKey = "actor_id"
)

// publicPaths are routes that don't require authentication.
var publicPaths = []string{
	"/api/v1/auth/",
	"/api/v1/docs",
}

func isPublicPath(path string) bool {
	for _, prefix := range publicPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// JWTAuth validates the Authorization header and injects tenant_id/actor_id into context.
// TODO(jwt-verify): Currently parses claims without signature verification.
// When the iam service shared secret is available, add HS256 validation.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse without verification for stub phase.
		parser := jwt.NewParser(jwt.WithoutClaimsValidation())
		token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		tenantID := claimToInt64(claims, "tenantId")
		actorID := claimToInt64(claims, "memberId")

		if tenantID == 0 {
			http.Error(w, `{"error":"missing tenant_id in token"}`, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyTenantID, tenantID)
		ctx = context.WithValue(ctx, ContextKeyActorID, actorID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func claimToInt64(claims jwt.MapClaims, key string) int64 {
	v, ok := claims[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	case int:
		return int64(n)
	default:
		return 0
	}
}

// TenantIDFromContext extracts tenant_id from the request context.
func TenantIDFromContext(ctx context.Context) int64 {
	v, _ := ctx.Value(ContextKeyTenantID).(int64)
	return v
}

// ActorIDFromContext extracts actor_id from the request context.
func ActorIDFromContext(ctx context.Context) int64 {
	v, _ := ctx.Value(ContextKeyActorID).(int64)
	return v
}
