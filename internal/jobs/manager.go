package jobs

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/isonlycoolie/grabb/internal/download"
)

// Runner downloads one link and returns the finished file.
type Runner func(ctx context.Context, link string, o download.Options, report func(download.Event)) (string, error)

// progressInterval limits how often progress alone is sent to subscribers.
const progressInterval = 250 * time.Millisecond

// Manager owns every job. All fields after mu are guarded by it. Changes
// are saved before subscribers hear about them.
type Manager struct {
	run     Runner
	store   store
	hub     hub
	workers int
	wake    chan struct{}

	mu      sync.Mutex
	jobs    map[string]*Job
	pending []string
	cancels map[string]context.CancelFunc
	sent    map[string]time.Time
}

// NewManager loads the library at path. Jobs that were still queued or
// running when grabb stopped are queued again; yt-dlp resumes partial files.
func NewManager(path string, workers int, run Runner) (*Manager, error) {
	m := &Manager{
		run: run, store: store{path}, workers: max(workers, 1), wake: make(chan struct{}, 1),
		jobs: map[string]*Job{}, cancels: map[string]context.CancelFunc{}, sent: map[string]time.Time{},
	}
	saved, err := m.store.load()
	if err != nil {
		return nil, err
	}
	slices.SortFunc(saved, func(a, b *Job) int { return a.CreatedAt.Compare(b.CreatedAt) })
	for _, j := range saved {
		m.jobs[j.ID] = j
		if j.active() {
			j.reset()
			m.pending = append(m.pending, j.ID)
		}
	}
	return m, nil
}

// Subscribe returns a channel of job changes and a function to stop them.
func (m *Manager) Subscribe() (<-chan Job, func()) { return m.hub.subscribe() }
