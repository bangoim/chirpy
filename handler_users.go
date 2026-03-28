package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bangoim/chirpy/internal/auth"
	"github.com/bangoim/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
