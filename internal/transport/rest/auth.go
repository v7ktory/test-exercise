package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		NotFoundErrorHandler(w, r)
		return
	}

	var input model.Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	if !model.IsEmailValid(input.Email) || input.Validate() != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	user := model.User{
		UUID:     uuid.New(),
		Email:    input.Email,
		Password: input.Password,
	}

	userID, err := h.Svc.SignUp(r.Context(), &user)
	if err != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userID); err != nil {
		InternalServerErrorHandler(w, r)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		NotFoundErrorHandler(w, r)
		return
	}

	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	var input model.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	if !model.IsEmailValid(input.Email) || input.Validate() != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	access, refresh, err := h.Svc.Login(r.Context(), userID, input.Email, input.Password)
	if err != nil {
		BadRequestErrorHandler(w, r)
		return
	}

	setRefreshTokenCookie(w, refresh.Token)

	response := model.AccessToken{
		ID:     access.ID,
		UserID: access.UserID,
		Token:  access.Token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerErrorHandler(w, r)
	}
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		NotFoundErrorHandler(w, r)
		return
	}
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		BadRequestErrorHandler(w, r)
		return
	}
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		NotFoundErrorHandler(w, r)
		return
	}

	access, refresh, err := h.Svc.Refresh(r.Context(), uuid.MustParse(userID), refreshCookie.Value)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	setRefreshTokenCookie(w, refresh.Token)

	response := model.AccessToken{
		ID:     access.ID,
		UserID: access.UserID,
		Token:  access.Token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerErrorHandler(w, r)
	}
}

// Устанавливаем refresh token в httpOnly куку
func setRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		Path:     "/auth",
		MaxAge:   3600 * 24 * 30, // 30 days
	}
	http.SetCookie(w, &cookie)
}
