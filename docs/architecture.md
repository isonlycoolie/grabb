# Architecture

grabb is one Go binary. It serves an embedded dashboard and a small JSON API, and drives
yt-dlp to do the downloading.

## Packages

| Package | Responsibility |
| --- | --- |
| `cmd/grabb` | flags, wiring and graceful shutdown |
| `internal/toolchain` | finds yt-dlp, ffmpeg and aria2c; turns yt-dlp errors into readable messages |
| `internal/media` | validates links and reads their details before downloading |
| `internal/download` | builds the yt-dlp command line, runs it and parses its progress |
| `internal/jobs` | queue, worker pool, live updates and the library file |
| `internal/httpapi` | HTTP routes, server-sent events and request guarding |
| `internal/desktop` | opens the browser and reveals files in the file manager |
| `web` | embeds `web/app`, the dashboard |

## A download, step by step

1. The dashboard posts the link to `POST /api/probe`. `media.Probe` runs `yt-dlp -J` and
   returns the title, thumbnail and available resolutions.
2. The person picks video or audio and a quality, and the dashboard posts to
   `POST /api/jobs`.
3. `jobs.Manager` saves the job to `library.json` and queues it. A free worker picks it up.
4. `download.Args` builds the command. Video takes the best video stream and the best
   audio stream, ranked by resolution, frame rate, HDR, codec and bitrate, and ffmpeg
   merges them. Audio extracts and converts to the chosen format.
5. `download.Run` reads yt-dlp's output. grabb asks yt-dlp to print marked lines for
   progress, stages, metadata and the final path, and also reads aria2c's progress lines.
6. Each event updates the job. The manager saves state changes first, then publishes them
   to every `GET /api/events` stream, and the dashboard updates the row in place.

## Dashboard

Plain ES modules with no build step, grouped by feature:

```
web/app/
  index.html          shell: sidebar, top bar, page outlet
  styles/             tokens, base, layout and one file per area
  scripts/
    main.js           routing and navigation
    api.js, store.js  server client and live job store
    lib/              DOM, icon, network mark, format and control helpers
    features/
      download/       link form, preview and quality choices
      jobs/           job rows and lists shared by pages
      library/        finished files with search and filters
      player/         in-app video and audio player
      settings/       configuration, engine versions and theme
```

The dashboard plays finished files from `GET /api/jobs/{id}/file?inline=1`, which supports
range requests so seeking works.

## Data

The library lives in the user config folder (`%AppData%\grabb\library.json` on Windows,
`~/.config/grabb/library.json` on Linux, `~/Library/Application Support/grabb` on macOS).
Jobs that were queued or running when grabb stopped are queued again on the next start,
and yt-dlp resumes their partial files.
