package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bangoim/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expiration := time.Hour
	if req.ExpiresInSeconds != nil && *req.ExpiresInSeconds > 0 && *req.ExpiresInSeconds < int(time.Hour.Seconds()) {
		expiration = time.Duration(*req.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, expiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	type loginResponse struct {
		User
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token: token,
	})
}
