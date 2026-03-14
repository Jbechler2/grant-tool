package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jbechler2/grant-tool/backend/internal/service"
)

type AuthHandler struct {
	authService  *service.AuthService
	isProduction bool
}

func NewAuthHandler(authService *service.AuthService, isProduction bool) *AuthHandler {
	return &AuthHandler{authService: authService, isProduction: isProduction}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	result, err := h.authService.Register(r.Context(), service.RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: r.UserAgent(),
		IpAddress: r.RemoteAddr,
	})
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email already in use")
			return
		}
		writeError(w, http.StatusInternalServerError, "registration failed")
		return
	}

	cookieDomain := ""
	if h.isProduction {
		cookieDomain = ".upwellgrantmanager.com"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    result.Token,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
		Path:     "/",
		Domain:   cookieDomain,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(result.RefreshExpiry).Seconds()),
		Path:     "/",
		Domain:   cookieDomain,
	})

	writeJSON(w, http.StatusCreated, authResponse{
		Email: result.User.Email,
		Role:  string(result.User.Role),
	})

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	result, err := h.authService.Login(r.Context(), service.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: r.UserAgent(),
		IpAddress: r.RemoteAddr,
	})

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if errors.Is(err, service.ErrTooManySessions) {
			writeError(w, http.StatusForbidden, "max sessions exceeded")
			return
		}
		writeError(w, http.StatusInternalServerError, "login failed")
		return
	}
	cookieDomain := ""
	if h.isProduction {
		cookieDomain = ".upwellgrantmanager.com"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    result.Token,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
		Path:     "/",
		Domain:   cookieDomain,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(result.RefreshExpiry).Seconds()),
		Path:     "/",
		Domain:   cookieDomain,
	})

	writeJSON(w, http.StatusOK, authResponse{
		Email: result.User.Email,
		Role:  string(result.User.Role),
	})

}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "refresh failed")
		return
	}

	result, err := h.authService.RotateToken(r.Context(), refreshToken.Value, service.RotateTokenInput{
		UserAgent: r.UserAgent(),
		IpAddress: r.RemoteAddr,
	})
	if err != nil {
		if errors.Is(err, service.ErrRefreshTokenNotFound) || errors.Is(err, service.ErrRefreshTokenExpired) {
			writeError(w, http.StatusUnauthorized, "refresh failed")
			return
		}
		writeError(w, http.StatusInternalServerError, "refresh failed")
		return
	}

	cookieDomain := ""
	if h.isProduction {
		cookieDomain = ".upwellgrantmanager.com"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    result.Token,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
		Path:     "/",
		Domain:   cookieDomain,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(result.RefreshExpiry).Seconds()),
		Path:     "/",
		Domain:   cookieDomain,
	})

	writeJSON(w, http.StatusOK, authResponse{
		Email: result.User.Email,
		Role:  string(result.User.Role),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil || refreshToken.Value == "" {
		h.clearCookies(w, h.isProduction)
		w.WriteHeader(http.StatusOK)
		return
	}

	if refreshToken != nil {
		err = h.authService.Logout(r.Context(), refreshToken.Value)
		if err != nil {
			log.Printf("logout cleanup failed: %v", err)
		}
	}

	h.clearCookies(w, h.isProduction)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) clearCookies(w http.ResponseWriter, isProduction bool) {
	cookieDomain := ""
	if h.isProduction {
		cookieDomain = ".upwellgrantmanager.com"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Path:     "/",
		Domain:   cookieDomain,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Path:     "/",
		Domain:   cookieDomain,
	})
}
