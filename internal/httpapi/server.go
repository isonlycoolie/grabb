// Package httpapi serves the dashboard and the JSON API behind it.
package httpapi

import (
	"context"
	"io/fs"
	"net/http"
	"strings"

	"github.com/isonlycoolie/grabb/internal/jobs"
	"github.com/isonlycoolie/grabb/internal/media"
)

// Prober reads the details of a link without downloading it.
type Prober func(ctx context.Context, link string) (media.Info, error)

// Revealer shows a file in the system file manager.
type Revealer func(path string) error

// Settings is what the settings page shows. It is fixed at startup.
type Settings struct {
	Version     string            `json:"version"`
	OutDir      string            `json:"outDir"`
	Workers     int               `json:"workers"`
	Connections int               `json:"connections"`
	Tools       map[string]string `json:"tools"`
}

type server struct {
	jobs     *jobs.Manager
	probe    Prober
	reveal   Revealer
	settings Settings
}

// New returns the handler for the API and the embedded dashboard.
func New(m *jobs.Manager, probe Prober, reveal Revealer, settings Settings, web fs.FS) http.Handler {
	s := &server{jobs: m, probe: probe, reveal: reveal, settings: settings}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/settings", s.getSettings)
	mux.HandleFunc("POST /api/probe", s.postProbe)
	mux.HandleFunc("GET /api/jobs", s.listJobs)
	mux.HandleFunc("POST /api/jobs", s.createJob)
	mux.HandleFunc("POST /api/jobs/{id}/cancel", s.cancelJob)
	mux.HandleFunc("POST /api/jobs/{id}/retry", s.retryJob)
	mux.HandleFunc("POST /api/jobs/{id}/reveal", s.revealJob)
	mux.HandleFunc("DELETE /api/jobs/{id}", s.deleteJob)
	mux.HandleFunc("GET /api/jobs/{id}/file", s.jobFile)
	mux.HandleFunc("GET /api/events", s.events)
	mux.Handle("GET /", http.FileServerFS(web))
	return guard(mux)
}

func (s *server) getSettings(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.settings)
}

func (s *server) postProbe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	if err := media.CheckURL(req.URL); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	info, err := s.probe(r.Context(), strings.TrimSpace(req.URL))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, info)
}
