package jobs

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/isonlycoolie/grabb/internal/download"
)

// execute runs one job to its end and records the outcome.
func (m *Manager) execute(parent context.Context, id string) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	m.mu.Lock()
	j := m.jobs[id]
	j.Status, j.Stage = StatusRunning, "Starting"
	link, opts := j.URL, j.Options
	m.cancels[id] = cancel
	m.publishLocked(j)
	m.mu.Unlock()

	file, err := m.run(ctx, link, opts, func(ev download.Event) { m.report(id, ev) })

	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cancels, id)
	delete(m.sent, id)
	j, ok := m.jobs[id]
	if !ok {
		return // removed while it ran
	}
	now := time.Now()
	j.Stage, j.FinishedAt = "", &now
	switch {
	case err == nil:
		j.Status, j.File = StatusDone, file
		if info, statErr := os.Stat(file); statErr == nil {
			j.Size = info.Size()
		}
	case parent.Err() != nil:
		// grabb is shutting down: leave it queued so the next start resumes it.
		j.Status, j.FinishedAt = StatusQueued, nil
	case ctx.Err() != nil:
		j.Status = StatusCanceled
	default:
		j.Status, j.Error = StatusFailed, err.Error()
	}
	m.saveLocked()
	m.publishLocked(j)
}

func (m *Manager) report(id string, ev download.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	j, ok := m.jobs[id]
	if !ok {
		return
	}
	stage := j.Stage
	j.apply(ev)
	onlyProgress := ev.Meta == nil && ev.File == "" && j.Stage == stage
	if onlyProgress && time.Since(m.sent[id]) < progressInterval {
		return
	}
	m.publishLocked(j)
}

func (m *Manager) publishLocked(j *Job) {
	m.sent[j.ID] = time.Now()
	m.hub.publish(*j)
}

func (m *Manager) saveLocked() {
	if err := m.store.save(m.listLocked()); err != nil {
		slog.Warn("could not save the library", "error", err)
	}
}
