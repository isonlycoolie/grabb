package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// store keeps jobs in one JSON file. Writes go to a temporary file first so
// a crash never leaves a half-written library.
type store struct {
	path string
}

func (s store) load() ([]*Job, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var jobs []*Job
	if err := json.Unmarshal(data, &jobs); err != nil {
		return nil, fmt.Errorf("read library %s: %w", s.path, err)
	}
	return jobs, nil
}

func (s store) save(jobs []Job) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// hub fans job updates out to subscribers. A subscriber that falls behind
// misses updates rather than blocking downloads; the dashboard reloads the
// full list whenever it reconnects.
type hub struct {
	mu   sync.Mutex
	subs map[chan Job]struct{}
}

func (h *hub) subscribe() (<-chan Job, func()) {
	ch := make(chan Job, 256)
	h.mu.Lock()
	if h.subs == nil {
		h.subs = make(map[chan Job]struct{})
	}
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch, func() {
		h.mu.Lock()
		delete(h.subs, ch)
		h.mu.Unlock()
	}
}

func (h *hub) publish(j Job) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs {
		select {
		case ch <- j:
		default:
		}
	}
}
