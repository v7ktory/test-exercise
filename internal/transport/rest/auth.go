package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
)

/*
Валидируем тело запроса и отправляем в сервисный слой
если всё ок возвращаем userID
*/
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

/*
Достаем userID из параметров запроса, email и password из тела запроса.
Передаем в сервисный слой и если всё ок создаем пару accessToken и refreshToken
AccessToken идет в header Authorization, refreshToken отправляем в куки
*/
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

	w.Header().Set("Authorization", "Bearer "+access.Token)
	setRefreshTokenCookie(w, refresh.Token)

	response := model.AccessToken{
		ID:     access.ID,
		UserID: access.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerErrorHandler(w, r)
	}
}

/*
Достаем userID из параметров запроса, refreshToken из куки и accessToken из header.
Передаем в сервисный слой и если всё ок обновляем пару
*/
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

	accessToken := extractAccessToken(r)
	if accessToken == "" {
		BadRequestErrorHandler(w, r)
		return
	}

	access, refresh, err := h.Svc.Refresh(r.Context(), uuid.MustParse(userID), accessToken, refreshCookie.Value)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	setRefreshTokenCookie(w, refresh.Token)

	response := model.AccessToken{
		ID:     access.ID,
		UserID: access.UserID,
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

// Извлекаем access token из header Authorization
func extractAccessToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}
