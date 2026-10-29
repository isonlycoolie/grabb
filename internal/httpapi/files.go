package httpapi

import (
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/isonlycoolie/grabb/internal/jobs"
)

func (s *server) revealJob(w http.ResponseWriter, r *http.Request) {
	job, ok := s.finished(w, r)
	if ok {
		respond(w, s.reveal(job.File))
	}
}

// jobFile streams a finished file: as an attachment to save a copy, or
// inline with ?inline=1 to play it in the browser.
func (s *server) jobFile(w http.ResponseWriter, r *http.Request) {
	job, ok := s.finished(w, r)
	if !ok {
		return
	}
	f, err := os.Open(job.File)
	if err != nil {
		writeError(w, http.StatusGone, errors.New("the file was moved or deleted"))
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	disposition := "attachment"
	if r.URL.Query().Get("inline") == "1" {
		disposition = "inline"
	}
	name := filepath.Base(job.File)
	w.Header().Set("Content-Disposition", mime.FormatMediaType(disposition, map[string]string{"filename": name}))
	http.ServeContent(w, r, name, info.ModTime(), f)
}

func (s *server) finished(w http.ResponseWriter, r *http.Request) (jobs.Job, bool) {
	job, ok := s.jobs.Get(r.PathValue("id"))
	if !ok || job.Status != jobs.StatusDone {
		writeError(w, http.StatusNotFound, jobs.ErrNotFound)
		return jobs.Job{}, false
	}
	return job, true
}
