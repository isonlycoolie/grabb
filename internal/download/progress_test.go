package download

import "testing"

func TestParseProgressTemplate(t *testing.T) {
	ev, ok := parseLine(markProgress + "1048576|NA|4194304.5|524288.0|6|none")
	if !ok || ev.Progress == nil {
		t.Fatalf("progress line not recognised: %+v", ev)
	}
	p := *ev.Progress
	if p.Downloaded != 1<<20 || p.Total != 4194304 || p.Speed != 524288 || p.ETA != 6 {
		t.Errorf("progress = %+v", p)
	}
	if ev.Stage != "Downloading audio" {
		t.Errorf("stage = %q, want Downloading audio", ev.Stage)
	}
}

func TestParseProgressUnknownValues(t *testing.T) {
	ev, _ := parseLine(markProgress + "NA|NA|NA|NA|NA|av01.0.13M.08")
	if ev.Progress.ETA != -1 || ev.Progress.Total != 0 || ev.Stage != "Downloading video" {
		t.Errorf("event = %+v / %+v", ev, ev.Progress)
	}
}

func TestParseStageMetaAndFile(t *testing.T) {
	if ev, ok := parseLine(markStage + "FFmpegMerger|started"); !ok || ev.Stage != "Merging" {
		t.Errorf("merge stage = %+v, %v", ev, ok)
	}
	if _, ok := parseLine(markStage + "FFmpegMerger|finished"); ok {
		t.Error("finished stages should be ignored")
	}
	ev, ok := parseLine(markMeta + `{"title":"Big Buck Bunny","extractor_key":"Youtube","duration":634.5}`)
	if !ok || ev.Meta.Title != "Big Buck Bunny" || ev.Meta.Site != "Youtube" {
		t.Errorf("meta = %+v, %v", ev.Meta, ok)
	}
	if ev, ok := parseLine(markFile + `C:\out\Youtube\a [x].mp4`); !ok || ev.File != `C:\out\Youtube\a [x].mp4` {
		t.Errorf("file = %+v, %v", ev, ok)
	}
	if _, ok := parseLine("[youtube] Extracting URL"); ok {
		t.Error("log lines should not parse as events")
	}
}
