package download

import (
	"slices"
	"testing"
)

var testConfig = Config{OutDir: "out", Connections: 16, FFmpeg: "ffmpeg", Aria2: "aria2c"}

// flag returns the value that follows name, or "" when name is absent.
func flag(args []string, name string) string {
	if i := slices.Index(args, name); i >= 0 && i+1 < len(args) {
		return args[i+1]
	}
	return ""
}

func TestArgsVideo(t *testing.T) {
	tests := []struct {
		name      string
		opts      Options
		sort      string
		container string
	}{
		{"best", Options{Mode: ModeVideo}, "res," + videoRank, "mp4/mkv"},
		{"capped", Options{Mode: ModeVideo, Height: 1080}, "res:1080," + videoRank, "mp4/mkv"},
		{"compatible", Options{Mode: ModeVideo, Compatible: true}, "vcodec:h264,acodec:aac,res," + videoRank, "mp4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := Args(testConfig, tt.opts, "https://example.com/v")
			if got := flag(args, "--format-sort"); got != tt.sort {
				t.Errorf("sort = %q, want %q", got, tt.sort)
			}
			if got := flag(args, "--merge-output-format"); got != tt.container {
				t.Errorf("container = %q, want %q", got, tt.container)
			}
		})
	}
}

func TestArgsAudio(t *testing.T) {
	args := Args(testConfig, Options{Mode: ModeAudio, AudioFormat: "flac"}, "https://example.com/a")
	if got := flag(args, "--audio-format"); got != "flac" {
		t.Errorf("audio format = %q, want flac", got)
	}
	if slices.Contains(args, "--format-sort") {
		t.Error("audio downloads must not sort video formats")
	}
}

func TestArgsEndWithSeparatedLink(t *testing.T) {
	args := Args(testConfig, Options{Mode: ModeVideo}, "https://example.com/v")
	if n := len(args); args[n-2] != "--" || args[n-1] != "https://example.com/v" {
		t.Errorf("args end with %q, want -- then the link", args[n-2:])
	}
}

func TestArgsForceUTF8(t *testing.T) {
	if got := flag(Args(testConfig, Options{Mode: ModeVideo}, "https://example.com/v"), "--encoding"); got != "utf-8" {
		t.Errorf("encoding = %q, want utf-8 so non-ASCII paths survive", got)
	}
}

func TestArgsAria2Optional(t *testing.T) {
	with := Args(testConfig, Options{Mode: ModeVideo}, "https://example.com/v")
	if flag(with, "--downloader") != "aria2c" {
		t.Error("aria2c should be used when it is installed")
	}
	cfg := testConfig
	cfg.Aria2 = ""
	if slices.Contains(Args(cfg, Options{Mode: ModeVideo}, "https://example.com/v"), "--downloader") {
		t.Error("no external downloader expected without aria2c")
	}
}

func TestNormalize(t *testing.T) {
	audio := Options{Mode: ModeAudio, Height: 720, Compatible: true}
	if err := audio.Normalize(); err != nil || audio.AudioFormat != "mp3" || audio.Height != 0 || audio.Compatible {
		t.Errorf("audio normalized to %+v, err %v", audio, err)
	}
	for _, bad := range []Options{{Mode: "gif"}, {Mode: ModeAudio, AudioFormat: "wav"}, {Height: 9000}} {
		if err := bad.Normalize(); err == nil {
			t.Errorf("Normalize(%+v) succeeded, want an error", bad)
		}
	}
}
