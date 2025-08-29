package api

import (
	"encoding/json"
	"net/http"

	"sledgehammer.echo-mesh.com/internal/clientbans"
)

type API struct {
	Store *clientbans.Store
}

func (a *API) ListBans(w http.ResponseWriter, r *http.Request) {
	bans, err := a.Store.List()
	if err != nil {
		http.Error(w, "failed to fetch bans", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bans)
}

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

	ban := clientbans.ClientBan{
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
