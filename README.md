<p align="center">
  <img src="web/app/assets/logo.svg" width="88" height="88" alt="grabb">
</p>

<h1 align="center">grabb</h1>

<p align="center">
  <strong>Paste a link. Get the file.</strong><br>
  A free, open source downloader for YouTube, TikTok, Instagram, X and 1,800 more sites,<br>
  in the best quality they serve. Private, fast, and it runs on your own computer.
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-3e63dd" alt="MIT license"></a>
  <img src="https://img.shields.io/badge/Windows%20%C2%B7%20macOS%20%C2%B7%20Linux-supported-3e63dd" alt="Windows, macOS and Linux">
  <img src="https://img.shields.io/badge/quality-up%20to%208K%2060fps-3e63dd" alt="Up to 8K 60fps">
</p>

<p align="center">
  <a href="https://github.com/isonlycoolie/grabb/raw/assets/grabb-readme/grabb-demo.mp4">
    <img src="https://github.com/isonlycoolie/grabb/raw/assets/grabb-readme/grabb-demo.gif" alt="grabb demo: paste a link, fetch it, choose audio, download and play it" width="100%">
  </a>
  <br>
  <sub>Paste a link, fetch it, choose a format, download, play. <a href="https://github.com/isonlycoolie/grabb/raw/assets/grabb-readme/grabb-demo.mp4">Watch in HD</a></sub>
</p>

## Get started

You need [Git](https://git-scm.com/downloads). Open a terminal and run the lines for your system.
The first run installs everything else grabb needs, which takes a few minutes. After that it
starts in seconds.

**Windows** (PowerShell or Command Prompt)

```powershell
git clone https://github.com/isonlycoolie/grabb.git; cd grabb; .\start.cmd
```

**macOS and Linux**

```sh
git clone https://github.com/isonlycoolie/grabb.git && cd grabb && ./start.sh
```

grabb opens in your browser at <http://127.0.0.1:7878>. Next time, open the `grabb` folder and run
`.\start.cmd` or `./start.sh` again.

> On macOS the script uses [Homebrew](https://brew.sh). On Linux it may ask for your password to
> install ffmpeg and aria2.

## How to use it

1. Copy a link from YouTube, TikTok, Instagram, X, Facebook or almost any video site.
2. Paste it into grabb. You see the title, the length and every quality on offer.
3. Choose **Video** or **Audio**, pick a quality, and press **Download**.

Files are saved to **Downloads/grabb**, one folder per site. Open the **Library** to play them
right inside grabb, show them in their folder, or save a copy.

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://github.com/isonlycoolie/grabb/raw/assets/grabb-readme/library-dark.png">
  <img src="https://github.com/isonlycoolie/grabb/raw/assets/grabb-readme/library-light.png" alt="The grabb library with finished downloads">
</picture>

## What you get

| | |
| --- | --- |
| **Best quality** | Video and sound are fetched separately at their best and joined into one file: 4K, 8K, 60fps and HDR when the site has them. |
| **Audio only** | Save songs and podcasts as MP3, M4A, Opus or FLAC, with the cover art inside. |
| **Plays everywhere** | One switch gives you an H.264 MP4 that opens on any phone, TV or video editor. |
| **Fast** | Many connections per file, so downloads go as fast as your internet allows. |
| **Built-in player** | Watch and listen inside grabb without hunting for the file. |
| **Private** | No account, no ads, no tracking. Nothing leaves your computer except the download itself. |
| **Light and dark** | Follows your system theme, or pick one in Settings. |

<img src="https://github.com/isonlycoolie/grabb/raw/assets/grabb-readme/player-dark.png" alt="Playing a song inside grabb">

## Options

Add any of these after `start.cmd` or `start.sh`, for example `./start.sh -out ~/Videos`.

| Option | What it does |
| --- | --- |
| `-out <folder>` | Save files somewhere other than Downloads/grabb. |
| `-addr 0.0.0.0:7878` | Use grabb from your phone on the same Wi-Fi. Only do this on networks you trust. |
| `-workers <n>` | How many downloads run at once (default 2). |
| `-no-open` | Do not open the browser on start. |

## Questions

**Is it really free?** Yes. grabb is MIT licensed and always will be.

**A site stopped working.** Sites change often, so update yt-dlp: `winget upgrade yt-dlp.yt-dlp`
on Windows, `brew upgrade yt-dlp` on macOS, or `yt-dlp -U` on Linux. Then try again. If it still fails, [open an issue](https://github.com/isonlycoolie/grabb/issues).

**Can I download anything I want?** Use grabb for your own content, content under an open
license, and content you have permission to save. Respect the rules of each site and the rights
of creators.

## For developers

grabb is a single Go program with the dashboard built in. It drives
[yt-dlp](https://github.com/yt-dlp/yt-dlp), [ffmpeg](https://ffmpeg.org), [aria2](https://aria2.github.io)
and [deno](https://deno.com). Read [docs/architecture.md](docs/architecture.md) for how it fits
together and [CONTRIBUTING.md](CONTRIBUTING.md) before you open a pull request. To report a
security issue, see [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE). Built by [isonlycoolie](https://github.com/isonlycoolie).

The dashboard ships the [Inter](https://rsms.me/inter/) typeface under the SIL Open Font License
and network marks from [Simple Icons](https://simpleicons.org) under CC0. Network names and marks
belong to their owners.
