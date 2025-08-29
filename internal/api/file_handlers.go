package api

import (
	"encoding/json"
	"net/http"

	"sledgehammer.echo-mesh.com/internal/models"
)

// RequestFileBan TODO Move into moderation package. Replace this with report instead.
func (a *API) RequestFileBan(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Filename    string `json:"filename"`
		Reason      string `json:"reason"`
		SubmittedBy string `json:"submittedBy"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	fileBan := models.FileBan{
		Filename:    input.Filename,
		Reason:      input.Reason,
		SubmittedBy: input.SubmittedBy,
	}

	if err := a.Store.RequestFileBan(fileBan); err != nil {
		http.Error(w, "failed to add ban", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fileBan)
}

func (a *API) FetchFileBanstatus(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
}

func (a *API) FetchFileBanInfo(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)

}
