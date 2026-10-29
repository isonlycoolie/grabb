package download

import "strconv"

// Config holds the settings shared by every download.
type Config struct {
	OutDir      string
	Connections int
	FFmpeg      string
	Aria2       string // empty leaves single-file streams to yt-dlp itself
}

const (
	outputTemplate   = "%(extractor_key)s/%(title).150B [%(id)s].%(ext)s"
	progressTemplate = "download:" + markProgress +
		"%(progress.downloaded_bytes)s|%(progress.total_bytes)s|%(progress.total_bytes_estimate)s|" +
		"%(progress.speed)s|%(progress.eta)s|%(info.vcodec)s"
	stageTemplate = "postprocess:" + markStage + "%(progress.postprocessor)s|%(progress.status)s"
	metaTemplate  = "before_dl:" + markMeta + "%(.{id,title,uploader,channel,extractor_key,thumbnail,duration})j"
	fileTemplate  = "after_move:" + markFile + "%(filepath)s"
	// Resolution first, then frame rate, HDR, codec efficiency and bitrate.
	videoRank = "fps,hdr:12,vcodec,br"
)

// Args builds the yt-dlp command line for one link.
func Args(cfg Config, o Options, link string) []string {
	args := []string{
		// Windows pipes default to the ANSI code page, which garbles paths.
		"--encoding", "utf-8",
		"--newline", "--no-colors", "--no-playlist", "--no-mtime",
		"--no-simulate", "--progress",
		"--ffmpeg-location", cfg.FFmpeg,
		"--paths", cfg.OutDir, "--output", outputTemplate,
		"--concurrent-fragments", strconv.Itoa(cfg.Connections),
		"--embed-metadata", "--embed-thumbnail", "--embed-chapters",
		"--progress-template", progressTemplate,
		"--progress-template", stageTemplate,
		"--print", metaTemplate,
		"--print", fileTemplate,
	}
	if cfg.Aria2 != "" {
		// Single-file streams get several connections each through aria2c;
		// DASH and HLS stay native and fetch fragments in parallel instead.
		n := strconv.Itoa(min(cfg.Connections, 16))
		args = append(args,
			"--downloader", cfg.Aria2,
			"--downloader", "dash,m3u8:native",
			"--downloader-args", "aria2c:-x"+n+" -s"+n+" -k1M --summary-interval=1 --console-log-level=warn")
	}

	if o.Mode == ModeAudio {
		return append(args, "--format", "ba/b", "--extract-audio",
			"--audio-format", o.AudioFormat, "--audio-quality", "0", "--", link)
	}
	res := "res"
	if o.Height > 0 {
		res += ":" + strconv.Itoa(o.Height)
	}
	sort, merge := res+","+videoRank, "mp4/mkv"
	if o.Compatible {
		sort, merge = "vcodec:h264,acodec:aac,"+sort, "mp4"
	}
	return append(args, "--format", "bv*+ba/b", "--format-sort", sort,
		"--merge-output-format", merge, "--", link)
}
