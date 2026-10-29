package jobs

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/isonlycoolie/grabb/internal/download"
)

// waitFor returns the first published update for id that has status.
func waitFor(t *testing.T, updates <-chan Job, id string, status Status) Job {
	t.Helper()
	timeout := time.After(5 * time.Second)
	for {
		select {
		case j := <-updates:
			if j.ID == id && j.Status == status {
				return j
			}
		case <-timeout:
			t.Fatalf("job %s never reached %s", id, status)
		}
	}
}

func newTestManager(t *testing.T, run Runner) (*Manager, <-chan Job, string) {
	t.Helper()
	dir := t.TempDir()
	m, err := NewManager(filepath.Join(dir, "library.json"), 1, run)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	updates, stop := m.Subscribe()
	t.Cleanup(stop)
	m.Start(ctx)
	return m, updates, dir
}

func TestJobFinishesAndPersists(t *testing.T) {
	file := filepath.Join(t.TempDir(), "clip.mp4")
	m, updates, dir := newTestManager(t, func(_ context.Context, _ string, _ download.Options, report func(download.Event)) (string, error) {
		report(download.Event{Meta: &download.Meta{Title: "Clip", Site: "Youtube"}})
		return file, os.WriteFile(file, []byte("12345"), 0o644)
	})
	job := m.Submit("https://example.com/v", download.Options{Mode: download.ModeVideo}, download.Meta{})
	done := waitFor(t, updates, job.ID, StatusDone)
	if done.Title != "Clip" || done.Size != 5 || done.File != file {
		t.Errorf("finished job = %+v", done)
	}

	reloaded, err := NewManager(filepath.Join(dir, "library.json"), 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := reloaded.Get(job.ID); !ok || got.Status != StatusDone {
		t.Errorf("reloaded job = %+v, found %v", got, ok)
	}
}

func TestJobFailureAndRetry(t *testing.T) {
	calls := 0
	m, updates, _ := newTestManager(t, func(context.Context, string, download.Options, func(download.Event)) (string, error) {
		calls++
		return "", errors.New("Video unavailable")
	})
	job := m.Submit("https://example.com/v", download.Options{}, download.Meta{})
	if failed := waitFor(t, updates, job.ID, StatusFailed); failed.Error != "Video unavailable" {
		t.Errorf("error = %q", failed.Error)
	}
	if err := m.Retry(job.ID); err != nil {
		t.Fatal(err)
	}
	waitFor(t, updates, job.ID, StatusFailed)
	if calls != 2 {
		t.Errorf("runner called %d times, want 2", calls)
	}
}

func TestCancelRunningJob(t *testing.T) {
	m, updates, _ := newTestManager(t, func(ctx context.Context, _ string, _ download.Options, _ func(download.Event)) (string, error) {
		<-ctx.Done()
		return "", ctx.Err()
	})
	job := m.Submit("https://example.com/v", download.Options{}, download.Meta{})
	waitFor(t, updates, job.ID, StatusRunning)
	if err := m.Cancel(job.ID); err != nil {
		t.Fatal(err)
	}
	waitFor(t, updates, job.ID, StatusCanceled)
	if err := m.Retry("missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("retry missing = %v, want ErrNotFound", err)
	}
}
