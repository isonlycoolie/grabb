// Package download turns a link and a set of choices into a finished file
// by driving yt-dlp and reading its progress.
package download

import (
	"errors"
	"fmt"
	"slices"
)

// Mode picks between a video file and an audio-only file.
type Mode string

const (
	ModeVideo Mode = "video"
	ModeAudio Mode = "audio"
)

// Options are the choices a person makes for one download.
type Options struct {
	Mode        Mode   `json:"mode"`
	Height      int    `json:"height"`      // video: cap the resolution, 0 for the best
	AudioFormat string `json:"audioFormat"` // audio: mp3, m4a, opus or flac
	Compatible  bool   `json:"compatible"`  // video: H.264 and AAC in MP4
}

var audioFormats = []string{"mp3", "m4a", "opus", "flac"}

// Normalize fills defaults, clears fields that do not apply to the mode and
// rejects values yt-dlp would not accept.
func (o *Options) Normalize() error {
	switch o.Mode {
	case "", ModeVideo:
		if o.Height < 0 || o.Height > 4320 {
			return errors.New("height must be between 0 and 4320")
		}
		o.Mode, o.AudioFormat = ModeVideo, ""
	case ModeAudio:
		if o.AudioFormat == "" {
			o.AudioFormat = "mp3"
		}
		if !slices.Contains(audioFormats, o.AudioFormat) {
			return fmt.Errorf("audio format %q is not supported", o.AudioFormat)
		}
		o.Height, o.Compatible = 0, false
	default:
		return fmt.Errorf("mode %q is not supported", o.Mode)
	}
	return nil
}
