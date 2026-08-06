package service

import (
	"blogapi/internal/db"
	"blogapi/internal/models"
	"blogapi/internal/repository"
	"blogapi/internal/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const refreshTokenTTL = 7 * 24 * time.Hour

// HandleRegister godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account with an email and password.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.RegisterRequest	true	"User registration data"
//	@Success		201		{object}	models.RegisterAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		409		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/users/register [post]
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	repo := repository.UserRepository{
		DB: db.DB,
	}
	var req models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid Credentials", nil)
		return
	}

	existingUser, err := repo.GetUserByEmail(req.Email)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "failed to signup", nil)
		return
	}

	if existingUser != nil {
		utils.JSON(w, http.StatusConflict, false, "Email already in use", nil)
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "SignUp failed, please try again later", nil)
		return
	}

	user, err := repo.CreateUser(req.Email, hashedPassword)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Could not register", nil)
		return
	}

	response := models.RegisterResponse{
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}

	utils.JSON(w, http.StatusCreated, true, "SignUp successfully", response)
}

// HandleLogin godoc
//
//	@Summary		Log in a user
//	@Description	Authenticates a user using an email and password.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.LoginRequest	true	"User login credentials"
//
// @Success 200 {object} models.LoginAPIResponse
// @Failure		400		{object}	utils.APIResponse
// @Failure		401		{object}	utils.APIResponse
// @Failure		500		{object}	utils.APIResponse
// @Router			/api/users/login [post]
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	userRepo := repository.UserRepository{
		DB: db.DB,
	}
	tokenRepo := repository.RefreshTokenRepository{
		DB: db.DB,
	}

	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "invalid credentials", nil)
		return
	}
	existingUser, err := userRepo.GetUserByEmail(req.Email)
	if existingUser == nil || err != nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid credentials", nil)
		return
	}
	err = utils.ChechHashedPassword(existingUser.Password, req.Password)
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid credentials", nil)
		return
	}

	accessToken, err := utils.GenerateJWT(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Login failed", nil)
		return
	}

	err = tokenRepo.DeleteUserRefreshTokens(existingUser.ID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Login failed", nil)
		return
	}

	refreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Login failed", nil)
		return
	}

	expiresAt := time.Now().Add(refreshTokenTTL).UTC()

	err = tokenRepo.CreateRefreshToken(existingUser.ID, refreshTokenStr, expiresAt)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Login failed", nil)
		return
	}

	response := models.LoginResponse{
		User: models.UserResponse{
			ID:        existingUser.ID,
			Email:     existingUser.Email,
			Role:      existingUser.Role,
			CreatedAt: existingUser.CreatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}

	utils.JSON(w, http.StatusOK, true, "Login successfully", response)
}

// HandleRefreshToken godoc
//
//	@Summary		Refresh authentication tokens
//	@Description	Exchanges a valid refresh token for a new access token and a new refresh token.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.RefreshTokenRequest	true	"Refresh token payload"
//	@Success		200		{object}	models.RefreshTokenAPIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		401		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/users/refresh [post]
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenRepo := repository.RefreshTokenRepository{
		DB: db.DB,
	}
	userRepo := repository.UserRepository{
		DB: db.DB,
	}

	var req models.RefreshTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Refresh token is required", nil)
		return
	}

	storedToken, err := tokenRepo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if storedToken == nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Invalid refresh token", nil)
		return
	}
	expiresAt, err := time.Parse(time.RFC3339, storedToken.ExpiresAt)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	if expiresAt.Before(time.Now()) {
		_ = tokenRepo.DeleteRefreshToken(req.RefreshToken)

		utils.JSON(w, http.StatusUnauthorized, false, "Refresh token expired", nil)
		return
	}

	user, err := userRepo.GetUserByID(storedToken.UserID)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}
	if user == nil {
		_ = tokenRepo.DeleteRefreshToken(req.RefreshToken)

		utils.JSON(w, http.StatusUnauthorized, false, "Invalid refresh token", nil)
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Token refresh failed", nil)
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Token refresh failed", nil)
		return
	}

	err = tokenRepo.DeleteUserRefreshTokens(user.ID)
	if err != nil && err != sql.ErrNoRows {
		utils.JSON(w, http.StatusInternalServerError, false, "Token refresh failed", nil)
		return
	}
	err = tokenRepo.CreateRefreshToken(
		user.ID,
		newRefreshToken,
		time.Now().Add(refreshTokenTTL).UTC(),
	)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Token refresh failed", nil)
		return
	}

	response := models.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}
	utils.JSON(w, http.StatusOK, true, "Token refreshed successfully", response)
}

// HandleLogout godoc
//
//	@Summary		Log out a user
//	@Description	Invalidates the provided refresh token and logs the user out.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.LogoutRequest	true	"Refresh token to invalidate"
//	@Success		200		{object}	utils.APIResponse
//	@Failure		400		{object}	utils.APIResponse
//	@Failure		401		{object}	utils.APIResponse
//	@Failure		500		{object}	utils.APIResponse
//	@Router			/api/users/logout [post]
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	tokenRepo := repository.RefreshTokenRepository{
		DB: db.DB,
	}

	var req models.LogoutRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		utils.JSON(w, http.StatusBadRequest, false, "Refresh token is required", nil)
		return
	}

	err = tokenRepo.DeleteRefreshToken(req.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusUnauthorized, false, "Invalid refresh token", nil)
			return
		}

		utils.JSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.JSON(w, http.StatusOK, true, "Logged out successfully", nil)
}
