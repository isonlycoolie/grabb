package download

import (
	"regexp"
	"time"
)

// aria2c prints readouts like "[#2089b0 12MiB/40MiB(30%) CN:16 DL:8.1MiB ETA:3s]"
// while it downloads single-file streams.
var aria2Readout = regexp.MustCompile(`\[#\w+ ([\d.]+)(\w+)/([\d.]+)(\w+)\(\d+%\)[^\]]*?DL:([\d.]+)(\w+)(?: ETA:(\w+))?`)

var units = map[string]float64{"B": 1, "KiB": 1 << 10, "MiB": 1 << 20, "GiB": 1 << 30, "TiB": 1 << 40}

func parseAria2(line string) (Event, bool) {
	m := aria2Readout.FindStringSubmatch(line)
	if m == nil {
		return Event{}, false
	}
	p := Progress{Downloaded: size(m[1], m[2]), Total: size(m[3], m[4]), Speed: float64(size(m[5], m[6])), ETA: -1}
	if d, err := time.ParseDuration(m[7]); err == nil {
		p.ETA = int(d.Seconds())
	}
	return Event{Stage: "Downloading", Progress: &p}, true
}

func size(value, unit string) int64 { return int64(number(value) * units[unit]) }
