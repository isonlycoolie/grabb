package download

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Markers prefix the lines grabb asks yt-dlp to print, so they can be told
// apart from ordinary log output.
const (
	markProgress = "grabb:progress:"
	markStage    = "grabb:stage:"
	markMeta     = "grabb:meta:"
	markFile     = "grabb:file:"
)

// Progress is the state of the stream being downloaded right now.
type Progress struct {
	Downloaded int64   `json:"downloaded"`
	Total      int64   `json:"total"`
	Speed      float64 `json:"speed"` // bytes per second
	ETA        int     `json:"eta"`   // seconds, -1 when unknown
}

// Meta is what yt-dlp knows about the media once extraction finishes.
type Meta struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Uploader  string  `json:"uploader"`
	Channel   string  `json:"channel"`
	Site      string  `json:"extractor_key"`
	Thumbnail string  `json:"thumbnail"`
	Duration  float64 `json:"duration"`
}

// Event is one thing yt-dlp reported. Stage may come with Progress; the
// other fields arrive alone.
type Event struct {
	Stage    string
	Progress *Progress
	Meta     *Meta
	File     string
}

func parseLine(line string) (Event, bool) {
	if rest, ok := strings.CutPrefix(line, markProgress); ok {
		return parseProgress(rest), true
	}
	if rest, ok := strings.CutPrefix(line, markStage); ok {
		name, status, _ := strings.Cut(rest, "|")
		return Event{Stage: stageName(name)}, status == "started"
	}
	if rest, ok := strings.CutPrefix(line, markMeta); ok {
		var m Meta
		return Event{Meta: &m}, json.Unmarshal([]byte(rest), &m) == nil
	}
	if rest, ok := strings.CutPrefix(line, markFile); ok {
		return Event{File: rest}, true
	}
	return Event{}, false
}

// parseProgress reads the fields of progressTemplate. yt-dlp prints NA for
// values it does not know yet.
func parseProgress(s string) Event {
	f := strings.Split(s, "|")
	for len(f) < 6 {
		f = append(f, "NA")
	}
	p := Progress{Downloaded: int64(number(f[0])), Total: int64(number(f[1])), Speed: number(f[3]), ETA: -1}
	if p.Total == 0 {
		p.Total = int64(number(f[2]))
	}
	if eta, err := strconv.ParseFloat(f[4], 64); err == nil {
		p.ETA = int(eta)
	}
	stage := "Downloading video"
	if f[5] == "none" {
		stage = "Downloading audio"
	}
	return Event{Stage: stage, Progress: &p}
}

func stageName(postprocessor string) string {
	switch {
	case strings.HasSuffix(postprocessor, "Merger"):
		return "Merging"
	case strings.HasSuffix(postprocessor, "ExtractAudio"):
		return "Converting audio"
	}
	return "Finishing"
}

func number(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
