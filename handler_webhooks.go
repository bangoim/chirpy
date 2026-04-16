package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := cfg.db.UpgradeUserToChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
