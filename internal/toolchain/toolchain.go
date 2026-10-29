// Package toolchain locates the programs grabb drives: yt-dlp to read
// sites, ffmpeg to merge and convert streams, and aria2c to open several
// connections per file.
package toolchain

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Tools holds absolute paths. Aria2 is optional and empty when missing.
type Tools struct {
	YtDlp  string
	FFmpeg string
	Aria2  string
}

// Locate finds every tool on PATH, falling back to the winget links folder
// so a fresh Windows install works before the shell restarts.
func Locate() (Tools, error) {
	t := Tools{YtDlp: find("yt-dlp"), FFmpeg: find("ffmpeg"), Aria2: find("aria2c")}
	if t.YtDlp == "" {
		return t, errors.New("yt-dlp not found: install it and make sure it is on PATH")
	}
	if t.FFmpeg == "" {
		return t, errors.New("ffmpeg not found: install it and make sure it is on PATH")
	}
	return t, nil
}

func find(name string) string {
	if p, err := exec.LookPath(name); err == nil {
		return p
	}
	if runtime.GOOS != "windows" {
		return ""
	}
	p := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "WinGet", "Links", name+".exe")
	if _, err := os.Stat(p); err != nil {
		return ""
	}
	return p
}

// Versions reports each tool's version for the settings page.
func (t Tools) Versions(ctx context.Context) map[string]string {
	return map[string]string{
		"yt-dlp": Version(ctx, t.YtDlp, "--version"),
		"ffmpeg": Version(ctx, t.FFmpeg, "-version"),
		"aria2c": Version(ctx, t.Aria2, "--version"),
	}
}

// Version returns the first line a tool prints for its version flag, or an
// empty string when the tool is missing or does not answer.
func Version(ctx context.Context, path, flag string) string {
	if path == "" {
		return ""
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, path, flag).Output()
	if err != nil {
		return ""
	}
	line, _, _ := strings.Cut(string(out), "\n")
	line, _, _ = strings.Cut(line, " Copyright")
	return strings.TrimSpace(line)
}

// Failure turns yt-dlp output into an error a person can act on. It keeps
// the last ERROR line and drops the extractor tag in front of it.
func Failure(output string, err error) error {
	msg := ""
	for _, line := range strings.Split(output, "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "ERROR: "); ok {
			msg = rest
		}
	}
	if msg == "" {
		return err
	}
	if strings.HasPrefix(msg, "[") {
		if _, rest, ok := strings.Cut(msg, ": "); ok {
			msg = rest
		}
	}
	return errors.New(msg)
}
