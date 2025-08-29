package api

import (
	"encoding/json"
	"net/http"

	"sledgehammer.echo-mesh.com/internal/models"
)

// AddClientBan TODO Move into moderation package. Replace this with report instead.
func (a *API) AddClientBan(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Character   string `json:"character"`
		World       string `json:"world"`
		LodestoneID string `json:"lodestoneId"`
		Reason      string `json:"reason"`
		SubmittedBy string `json:"submittedBy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	ban := models.ClientBan{
		Character:   input.Character,
		World:       input.World,
		LodestoneID: input.LodestoneID,
		Reason:      input.Reason,
		SubmittedBy: input.SubmittedBy,
	}

	if err := a.Store.AddClientBan(ban); err != nil {
		http.Error(w, "failed to add ban", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ban)
}

func (a *API) FetchPlayerBanStatus(writer http.ResponseWriter, request *http.Request) {
	a.Logger.Debugf("checking ban status for %s", request.URL.Path)
	params, ok := validatePlayerParams(writer, request)
	if !ok {
		a.Logger.Errorf("invalid input: character ID and world are required")
		return // error response already written
	}

	exists, err := a.Store.DoesClientBanExist(params.CharacterID, params.World) // /bans/player/{id}?world=...
	if err != nil {
		a.Logger.Errorf("failed to check ban status for character %s on %s", params.CharacterID, params.World)
		http.Error(writer, "failed to check ban status", http.StatusInternalServerError)
		return
	}

	if exists {
		a.Logger.Infof("ban found for character %s on %s", params.CharacterID, params.World)
		writer.WriteHeader(http.StatusOK)
	} else {
		a.Logger.Infof("no ban found for character %s on %s", params.CharacterID, params.World)
		writer.WriteHeader(http.StatusNotFound)
	}
}

func (a *API) FetchPlayerBanInfo(writer http.ResponseWriter, request *http.Request) {
	a.Logger.Debugf("fetching ban info for %s", request.URL.Path)
	params, ok := validatePlayerParams(writer, request)
	if !ok {
		a.Logger.Errorf("invalid input: character ID and world are required")
		return // error response already written
	}
	banInfo, err := a.Store.GetPlayerBanInfo(params.CharacterID, params.World)
	if err != nil {
		a.Logger.Errorf("failed to get ban info for character %s on %s", params.CharacterID, params.World)
		http.Error(writer, "failed to get ban info", http.StatusInternalServerError)
		return
	}
	a.Logger.Infof("fetched ban info for character %s on %s", params.CharacterID, params.World)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(banInfo)
}
