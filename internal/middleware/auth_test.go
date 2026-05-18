package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildTestJWT creates an unsigned JWT with the given claims (header.payload.signature).
func buildTestJWT(claims map[string]interface{}) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, _ := json.Marshal(claims)
	payloadEnc := base64.RawURLEncoding.EncodeToString(payload)
	sig := base64.RawURLEncoding.EncodeToString([]byte("stub-signature"))
	return fmt.Sprintf("%s.%s.%s", header, payloadEnc, sig)
}

func TestJWTAuth_PublicPath_NoToken(t *testing.T) {
	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/v1/auth/phone", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestJWTAuth_PublicPath_AuthCode(t *testing.T) {
	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/v1/auth/code", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestJWTAuth_PublicPath_Refresh(t *testing.T) {
	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestJWTAuth_MissingAuthHeader(t *testing.T) {
	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "missing authorization header")
}

func TestJWTAuth_InvalidFormat_NoBearerPrefix(t *testing.T) {
	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Token abc123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid authorization header format")
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTAuth_MissingTenantID(t *testing.T) {
	token := buildTestJWT(map[string]interface{}{
		"memberId": 42,
	})

	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "missing tenant_id")
}

func TestJWTAuth_ValidToken_ExtractsClaims(t *testing.T) {
	token := buildTestJWT(map[string]interface{}{
		"tenantId": float64(10),
		"memberId": float64(42),
	})

	var gotTenantID, gotActorID int64

	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTenantID = TenantIDFromContext(r.Context())
		gotActorID = ActorIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), gotTenantID)
	assert.Equal(t, int64(42), gotActorID)
}

func TestJWTAuth_BearerCaseInsensitive(t *testing.T) {
	token := buildTestJWT(map[string]interface{}{
		"tenantId": float64(5),
		"memberId": float64(1),
	})

	handler := JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "BEARER "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTenantIDFromContext_NoValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	assert.Equal(t, int64(0), TenantIDFromContext(req.Context()))
}

func TestActorIDFromContext_NoValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	assert.Equal(t, int64(0), ActorIDFromContext(req.Context()))
}
