package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/isonlycoolie/grabb/internal/download"
	"github.com/isonlycoolie/grabb/internal/jobs"
	"github.com/isonlycoolie/grabb/internal/media"
)

type createRequest struct {
	URL string `json:"url"`
	download.Options
	Known download.Meta `json:"known"` // details the dashboard already read
}

func (s *server) listJobs(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.jobs.List())
}

func (s *server) createJob(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if !readJSON(w, r, &req) {
		return
	}
	if err := media.CheckURL(req.URL); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := req.Options.Normalize(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, s.jobs.Submit(strings.TrimSpace(req.URL), req.Options, req.Known))
}

func (s *server) cancelJob(w http.ResponseWriter, r *http.Request) {
	respond(w, s.jobs.Cancel(r.PathValue("id")))
}

func (s *server) retryJob(w http.ResponseWriter, r *http.Request) {
	respond(w, s.jobs.Retry(r.PathValue("id")))
}

func (s *server) deleteJob(w http.ResponseWriter, r *http.Request) {
	respond(w, s.jobs.Remove(r.PathValue("id")))
}

func respond(w http.ResponseWriter, err error) {
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, jobs.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	default:
		writeError(w, http.StatusConflict, err)
	}
}
