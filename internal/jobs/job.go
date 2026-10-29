// Package jobs queues downloads, runs a few at a time, tells subscribers
// about every change and keeps the library of finished files on disk.
package jobs

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/isonlycoolie/grabb/internal/download"
)

// Status is where a job is in its life.
type Status string

const (
	StatusQueued   Status = "queued"
	StatusRunning  Status = "running"
	StatusDone     Status = "done"
	StatusFailed   Status = "failed"
	StatusCanceled Status = "canceled"
	// StatusRemoved is only sent to subscribers when a job is deleted.
	StatusRemoved Status = "removed"
)

// Job is one link being downloaded, or one finished file in the library.
type Job struct {
	ID         string            `json:"id"`
	URL        string            `json:"url"`
	Options    download.Options  `json:"options"`
	Status     Status            `json:"status"`
	Stage      string            `json:"stage,omitempty"`
	Title      string            `json:"title,omitempty"`
	Uploader   string            `json:"uploader,omitempty"`
	Site       string            `json:"site,omitempty"`
	Thumbnail  string            `json:"thumbnail,omitempty"`
	Duration   float64           `json:"duration,omitempty"`
	Progress   download.Progress `json:"progress"`
	File       string            `json:"file,omitempty"`
	Size       int64             `json:"size,omitempty"`
	Error      string            `json:"error,omitempty"`
	CreatedAt  time.Time         `json:"createdAt"`
	FinishedAt *time.Time        `json:"finishedAt,omitempty"`
}

func (j *Job) active() bool { return j.Status == StatusQueued || j.Status == StatusRunning }

func (j *Job) apply(ev download.Event) {
	if ev.Stage != "" {
		j.Stage = ev.Stage
	}
	if ev.Progress != nil {
		j.Progress = *ev.Progress
	}
	if m := ev.Meta; m != nil {
		j.Title = cmp.Or(m.Title, j.Title)
		j.Uploader = cmp.Or(m.Uploader, m.Channel, j.Uploader)
		j.Site = cmp.Or(m.Site, j.Site)
		j.Thumbnail = cmp.Or(m.Thumbnail, j.Thumbnail)
		j.Duration = cmp.Or(m.Duration, j.Duration)
	}
	if ev.File != "" {
		j.File = ev.File
	}
}

// reset clears the outcome of a previous run so the job can be queued again.
func (j *Job) reset() {
	j.Status, j.Stage, j.Error = StatusQueued, "", ""
	j.Progress, j.FinishedAt = download.Progress{}, nil
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
