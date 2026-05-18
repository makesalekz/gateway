package clients

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthByPhoneRequest struct {
	Phone string `json:"phone"`
}

type AuthByEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthByCodeRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type IAMClient interface {
	AuthByPhone(ctx context.Context, req *AuthByPhoneRequest) (*AuthResponse, error)
	AuthByEmail(ctx context.Context, req *AuthByEmailRequest) (*AuthResponse, error)
	AuthByCode(ctx context.Context, req *AuthByCodeRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AuthResponse, error)
}

type iamStub struct{}

func NewIAMClient() IAMClient {
	return &iamStub{}
}

// stubSigningKey is used only by the stub IAM client to produce valid JWTs for development.
var stubSigningKey = []byte("stub-secret-key-for-development-only")

func generateStubToken() string {
	claims := jwt.MapClaims{
		"tenantId": 10,
		"memberId": 5,
		"iss":      "iam-stub",
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString(stubSigningKey)
	return s
}

func generateStubRefreshToken() string {
	claims := jwt.MapClaims{
		"tenantId": 10,
		"memberId": 5,
		"iss":      "iam-stub",
		"type":     "refresh",
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString(stubSigningKey)
	return s
}

func (s *iamStub) AuthByPhone(_ context.Context, _ *AuthByPhoneRequest) (*AuthResponse, error) {
	return &AuthResponse{AccessToken: generateStubToken(), RefreshToken: generateStubRefreshToken(), ExpiresIn: 3600}, nil
}

func (s *iamStub) AuthByEmail(_ context.Context, _ *AuthByEmailRequest) (*AuthResponse, error) {
	return &AuthResponse{AccessToken: generateStubToken(), RefreshToken: generateStubRefreshToken(), ExpiresIn: 3600}, nil
}

func (s *iamStub) AuthByCode(_ context.Context, _ *AuthByCodeRequest) (*AuthResponse, error) {
	return &AuthResponse{AccessToken: generateStubToken(), RefreshToken: generateStubRefreshToken(), ExpiresIn: 3600}, nil
}

func (s *iamStub) RefreshToken(_ context.Context, _ *RefreshTokenRequest) (*AuthResponse, error) {
	return &AuthResponse{AccessToken: generateStubToken(), RefreshToken: generateStubRefreshToken(), ExpiresIn: 3600}, nil
}
