<p align="center">
  <img src="web/app/assets/logo.svg" width="72" height="72" alt="grabb logo">
</p>

<h1 align="center">grabb</h1>

<p align="center">
  Free, open source downloader for social networks. Paste a link, get the video or audio
  in the best quality the site serves. Runs on your own machine.
</p>

## Why grabb

- **Real top quality.** Picks the best video and the best audio separately and merges them
  with ffmpeg, so you get 4K, 8K, 60fps and HDR when the site has them.
- **Fast.** Fetches 16 fragments at once and opens 16 connections per file through aria2c.
- **Private.** No accounts, no tracking, no ads. Files go straight to your disk.
- **One binary.** A Go server with the dashboard embedded. No Node, no database, no runtime
  dependencies beyond the tools that do the downloading.
- **Many sites.** YouTube, TikTok, Instagram, X, Facebook, Vimeo, SoundCloud, Reddit and
  1,800 more through [yt-dlp](https://github.com/yt-dlp/yt-dlp).

## Requirements

| Tool | Why | Required |
| --- | --- | --- |
| [Go 1.25+](https://go.dev/dl/) | builds grabb | yes |
| [yt-dlp](https://github.com/yt-dlp/yt-dlp) | reads sites and streams | yes |
| [ffmpeg](https://ffmpeg.org/) | merges and converts | yes |
| [deno](https://deno.com/) | lets yt-dlp read YouTube | for YouTube |
| [aria2](https://aria2.github.io/) | parallel connections | recommended |

```sh
# Windows
winget install yt-dlp.yt-dlp Gyan.FFmpeg aria2.aria2 DenoLand.Deno
# macOS
brew install yt-dlp ffmpeg aria2 deno
# Debian or Ubuntu (install yt-dlp and deno from their releases for the latest versions)
sudo apt install ffmpeg aria2
```

## Run

```sh
git clone https://github.com/isonlycoolie/grabb.git
cd grabb
go run ./cmd/grabb
```

The dashboard opens at <http://127.0.0.1:7878>. Finished files land in `Downloads/grabb`,
one folder per site.

| Flag | Default | Meaning |
| --- | --- | --- |
| `-addr` | `127.0.0.1:7878` | where the dashboard listens; use `0.0.0.0:7878` to reach it from your phone |
| `-out` | `~/Downloads/grabb` | folder for finished files |
| `-workers` | `2` | downloads that run at the same time |
| `-connections` | `16` | parallel connections per file |
| `-no-open` | `false` | do not open a browser on start |

Build a single binary with `go build -o bin/grabb ./cmd/grabb`.

## How it works

```
browser ──HTTP/SSE──> httpapi ──> jobs (queue, workers, library.json)
                         │              │
                         └──> media     └──> download ──> yt-dlp ──> ffmpeg / aria2c
```

See [docs/architecture.md](docs/architecture.md) for the packages and the request flow.

## Use it responsibly

grabb is for content you own, content under an open license, and content you have
permission to save. Respect the terms of the sites you use and the rights of creators.

## Contributing

Issues and pull requests are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) for the
branch model, commit rules and code style, and [SECURITY.md](SECURITY.md) to report a
vulnerability privately.

## License

[MIT](LICENSE). Built by [isonlycoolie](https://github.com/isonlycoolie).

The dashboard ships the [Inter](https://rsms.me/inter/) typeface under the SIL Open Font
License and network marks from [Simple Icons](https://simpleicons.org/) under CC0. Network
names and marks belong to their owners.
