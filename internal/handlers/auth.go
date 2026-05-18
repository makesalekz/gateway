package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/makesalekz/gateway/internal/clients"
)

type AuthHandler struct {
	iam clients.IAMClient
}

func NewAuthHandler(iam clients.IAMClient) *AuthHandler {
	return &AuthHandler{iam: iam}
}

func (h *AuthHandler) AuthByPhone(w http.ResponseWriter, r *http.Request) {
	var req clients.AuthByPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Phone == "" {
		writeError(w, http.StatusBadRequest, "phone is required")
		return
	}

	resp, err := h.iam.AuthByPhone(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "auth failed")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) AuthByEmail(w http.ResponseWriter, r *http.Request) {
	var req clients.AuthByEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}

	resp, err := h.iam.AuthByEmail(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "auth failed")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) AuthByCode(w http.ResponseWriter, r *http.Request) {
	var req clients.AuthByCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Phone == "" || req.Code == "" {
		writeError(w, http.StatusBadRequest, "phone and code are required")
		return
	}

	resp, err := h.iam.AuthByCode(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "auth failed")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req clients.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	resp, err := h.iam.RefreshToken(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "refresh failed")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
