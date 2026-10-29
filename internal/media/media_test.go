package media

import (
	"slices"
	"testing"
)

func TestCheckURL(t *testing.T) {
	for link, ok := range map[string]bool{
		"https://www.youtube.com/watch?v=aqz-KE-bpKQ": true,
		"  http://vm.tiktok.com/abc  ":                true,
		"--exec calc":                                 false,
		"file:///etc/passwd":                          false,
		"youtube.com/watch?v=1":                       false,
		"https://":                                    false,
	} {
		if err := CheckURL(link); (err == nil) != ok {
			t.Errorf("CheckURL(%q) error = %v, want ok %v", link, err, ok)
		}
	}
}

func TestInfoResolutions(t *testing.T) {
	raw := rawInfo{Uploader: "", Channel: "Blender", ExtractorKey: "Youtube"}
	raw.Formats = []format{
		{Width: 1920, Height: 1080, VCodec: "avc1"},
		{Width: 3840, Height: 2160, VCodec: "av01"},
		{Width: 1080, Height: 1920, VCodec: "h264"}, // portrait 1080p
		{Width: 1920, Height: 818, VCodec: "vp9"},   // ultra-wide 1080p
		{Width: 1280, Height: 544, VCodec: "vp9"},   // ultra-wide 720p
		{VCodec: "none"},
	}
	info := raw.info()
	if want := []int{2160, 1080, 720}; !slices.Equal(info.Resolutions, want) {
		t.Errorf("resolutions = %v, want %v", info.Resolutions, want)
	}
	if info.Uploader != "Blender" || info.Site != "Youtube" {
		t.Errorf("uploader/site = %q/%q", info.Uploader, info.Site)
	}
}
