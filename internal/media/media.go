// Package media reads what a link offers before anything is downloaded.
package media

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/isonlycoolie/grabb/internal/toolchain"
)

// Info describes one video or post.
type Info struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Uploader    string  `json:"uploader"`
	Site        string  `json:"site"`
	Thumbnail   string  `json:"thumbnail"`
	Duration    float64 `json:"duration"`
	Resolutions []int   `json:"resolutions"` // shorter side of each video stream, highest first
}

type rawInfo struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Uploader     string   `json:"uploader"`
	Channel      string   `json:"channel"`
	ExtractorKey string   `json:"extractor_key"`
	Thumbnail    string   `json:"thumbnail"`
	Duration     float64  `json:"duration"`
	Formats      []format `json:"formats"`
}

type format struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	VCodec string `json:"vcodec"`
}

// CheckURL accepts absolute http and https links only, so nothing that
// looks like a flag or a local path reaches yt-dlp.
func CheckURL(raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("paste a full link that starts with http:// or https://")
	}
	return nil
}

// Probe asks yt-dlp for the details of link without downloading it.
func Probe(ctx context.Context, ytdlp, link string) (Info, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, ytdlp, "-J", "--encoding", "utf-8", "--no-playlist", "--no-warnings", "--", link)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return Info{}, toolchain.Failure(stderr.String(), err)
	}
	var raw rawInfo
	if err := json.Unmarshal(out, &raw); err != nil {
		return Info{}, fmt.Errorf("read link details: %w", err)
	}
	return raw.info(), nil
}

func (r rawInfo) info() Info {
	info := Info{
		ID: r.ID, Title: r.Title, Uploader: cmp.Or(r.Uploader, r.Channel),
		Site: r.ExtractorKey, Thumbnail: r.Thumbnail, Duration: r.Duration,
	}
	for _, f := range r.Formats {
		res := f.Height
		if f.Width > 0 && f.Width < res {
			res = f.Width // portrait video: quality follows the shorter side
		}
		if res > 0 && f.VCodec != "none" && !slices.Contains(info.Resolutions, res) {
			info.Resolutions = append(info.Resolutions, res)
		}
	}
	slices.SortFunc(info.Resolutions, func(a, b int) int { return b - a })
	return info
}
