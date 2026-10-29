package download

import "testing"

func TestParseAria2Readout(t *testing.T) {
	ev, ok := parseLine("[#2089b0 12MiB/40MiB(30%) CN:16 DL:8.0MiB ETA:1m3s]")
	if !ok {
		t.Fatal("aria2 readout not recognised")
	}
	p := *ev.Progress
	if p.Downloaded != 12<<20 || p.Total != 40<<20 || p.Speed != 8<<20 || p.ETA != 63 {
		t.Errorf("progress = %+v", p)
	}
	if _, ok := parseLine("[#2089b0 0B/0B CN:1 DL:0B]"); ok {
		t.Error("readout without a percentage should be ignored")
	}
}
