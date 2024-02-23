package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/v7ktory/test/internal/model"
	"github.com/v7ktory/test/pkg/validation"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input model.SignUpInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	if !validation.IsEmailValid(input.Email) {
		BadRequestErrorHandler(w, r)
		return
	}

	user := model.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     input.Password,
		RegisteredAt: time.Now(),
	}

	access, refresh, err := h.Svc.Auth.SignUp(r.Context(), &user)
	if err != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	setRefreshTokenCookie(w, refresh.Token)

	response := model.Response{
		AccessToken: model.AccessToken{
			Token:  access.Token,
			ID:     access.ID,
			UserID: access.UserID,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	if !validation.IsEmailValid(input.Email) {
		BadRequestErrorHandler(w, r)
		return
	}

	access, refresh, err := h.Svc.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	setRefreshTokenCookie(w, refresh.Token)

	response := model.Response{
		AccessToken: model.AccessToken{
			Token:  access.Token,
			ID:     access.ID,
			UserID: access.UserID,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func setRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		Path:     "/api/auth",
		MaxAge:   3600 * 24 * 30, // 30 days
	}
	http.SetCookie(w, &cookie)
}
