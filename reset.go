package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := cfg.db.DeleteAllUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete users")
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and all users deleted"))
}
