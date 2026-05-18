package clients

import (
	"context"
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

func (s *iamStub) AuthByPhone(_ context.Context, _ *AuthByPhoneRequest) (*AuthResponse, error) {
	return &AuthResponse{AccessToken: "stub-token", RefreshToken: "stub-refresh", ExpiresIn: 3600}, nil
}

func (s *iamStub) AuthByEmail(_ context.Context, _ *AuthByEmailRequest) (*AuthResponse, error) {
	return &AuthResponse{AccessToken: "stub-token", RefreshToken: "stub-refresh", ExpiresIn: 3600}, nil
}

func (s *iamStub) AuthByCode(_ context.Context, _ *AuthByCodeRequest) (*AuthResponse, error) {
	return &AuthResponse{AccessToken: "stub-token", RefreshToken: "stub-refresh", ExpiresIn: 3600}, nil
}

func (s *iamStub) RefreshToken(_ context.Context, _ *RefreshTokenRequest) (*AuthResponse, error) {
	return &AuthResponse{AccessToken: "stub-token", RefreshToken: "stub-refresh", ExpiresIn: 3600}, nil
}
