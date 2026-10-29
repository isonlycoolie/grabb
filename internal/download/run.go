package download

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/isonlycoolie/grabb/internal/toolchain"
)

// Run starts yt-dlp with args and calls report for every event it prints.
// It returns the path of the finished file.
func Run(ctx context.Context, ytdlp string, args []string, report func(Event)) (string, error) {
	cmd := exec.CommandContext(ctx, ytdlp, args...)
	cmd.WaitDelay = 5 * time.Second
	out, in := io.Pipe()
	cmd.Stdout, cmd.Stderr = in, in
	if err := cmd.Start(); err != nil {
		return "", err
	}
	go func() { in.CloseWithError(cmd.Wait()) }()

	var file string
	var tail []string // recent log lines, for the error message
	sc := bufio.NewScanner(out)
	sc.Buffer(make([]byte, 0, 64<<10), 1<<20)
	sc.Split(scanLines)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		ev, ok := parseLine(line)
		if !ok {
			if line != "" {
				tail = append(tail[max(0, len(tail)-19):], line)
			}
			continue
		}
		if ev.File != "" {
			file = ev.File
		}
		report(ev)
	}
	err := sc.Err()
	_, _ = io.Copy(io.Discard, out) // let yt-dlp exit if scanning stopped early

	switch {
	case ctx.Err() != nil:
		return "", ctx.Err()
	case err != nil:
		return "", toolchain.Failure(strings.Join(tail, "\n"), err)
	case file == "":
		return "", errors.New("yt-dlp finished without saving a file")
	}
	return file, nil
}

// scanLines splits on newlines and on the bare carriage returns that
// progress readouts use to redraw themselves.
func scanLines(data []byte, atEOF bool) (int, []byte, error) {
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		return i + 1, data[:i], nil
	}
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}
