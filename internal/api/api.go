package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"sledgehammer.echo-mesh.com/internal/database"
)

type API struct {
	Store  *database.Store
	Logger *zap.SugaredLogger
}

type PlayerParams struct {
	CharacterID string
	World       string
}

func validatePlayerParams(w http.ResponseWriter, r *http.Request) (*PlayerParams, bool) {
	characterID := chi.URLParam(r, "id")
	world := chi.URLParam(r, "world")

	// Basic validation
	if characterID == "" || world == "" {
		http.Error(w, "invalid input: character ID and world are required", http.StatusBadRequest)
		return nil, false
	}

	// Trim whitespace
	characterID = strings.TrimSpace(characterID)
	world = strings.TrimSpace(world)

	// TODO: Add check against list of current existing FFXIV worlds/servers

	return &PlayerParams{
		CharacterID: characterID,
		World:       world,
	}, true
}
