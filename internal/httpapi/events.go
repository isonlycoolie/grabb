package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// events streams job changes as server-sent events until the client leaves.
func (s *server) events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming is not supported", http.StatusInternalServerError)
		return
	}
	updates, stop := s.jobs.Subscribe()
	defer stop()

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	fmt.Fprint(w, "retry: 2000\n\n")
	flusher.Flush()

	ping := time.NewTicker(20 * time.Second)
	defer ping.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ping.C:
			fmt.Fprint(w, ": ping\n\n")
		case job := <-updates:
			data, err := json.Marshal(job)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: job\ndata: %s\n\n", data)
		}
		flusher.Flush()
	}
}
