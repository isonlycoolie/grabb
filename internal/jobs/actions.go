package jobs

import (
	"errors"
	"time"

	"github.com/isonlycoolie/grabb/internal/download"
)

var (
	ErrNotFound     = errors.New("download not found")
	ErrNotRetryable = errors.New("only failed or canceled downloads can be retried")
)

// Submit queues a download. known carries details the dashboard already
// read from the link so the job shows a title before yt-dlp starts.
func (m *Manager) Submit(link string, o download.Options, known download.Meta) Job {
	m.mu.Lock()
	defer m.mu.Unlock()
	j := &Job{ID: newID(), URL: link, Options: o, Status: StatusQueued, CreatedAt: time.Now()}
	j.apply(download.Event{Meta: &known})
	m.jobs[j.ID] = j
	m.enqueueLocked(j)
	return *j
}

// Cancel stops a queued or running job. Finished jobs are left alone.
func (m *Manager) Cancel(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	j, ok := m.jobs[id]
	switch {
	case !ok:
		return ErrNotFound
	case j.Status == StatusQueued:
		j.Status = StatusCanceled
		m.saveLocked()
		m.publishLocked(j)
	case j.Status == StatusRunning:
		m.cancels[id]()
	}
	return nil
}

// Retry queues a failed or canceled job again.
func (m *Manager) Retry(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	j, ok := m.jobs[id]
	switch {
	case !ok:
		return ErrNotFound
	case j.Status != StatusFailed && j.Status != StatusCanceled:
		return ErrNotRetryable
	}
	j.reset()
	m.enqueueLocked(j)
	return nil
}

// Remove deletes a job, stopping it first if needed. A downloaded file
// stays in the output folder.
func (m *Manager) Remove(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.jobs[id]; !ok {
		return ErrNotFound
	}
	if cancel := m.cancels[id]; cancel != nil {
		cancel()
	}
	delete(m.jobs, id)
	m.saveLocked()
	m.hub.publish(Job{ID: id, Status: StatusRemoved})
	return nil
}

func (m *Manager) enqueueLocked(j *Job) {
	m.pending = append(m.pending, j.ID)
	m.saveLocked()
	m.publishLocked(j)
	m.signal()
}
