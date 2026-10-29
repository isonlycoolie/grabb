# Contributing to grabb

Thanks for helping. This guide keeps the history clean and reviews quick.

## Branches

| Branch | Purpose |
| --- | --- |
| `main` | released code. Protected; changes arrive only through reviewed pull requests from `develop`. |
| `develop` | integration branch. Feature branches start here and merge back here. |
| `feat/<topic>` | a new capability, for example `feat/playlist-downloads` |
| `fix/<topic>` | a bug fix, for example `fix/tiktok-thumbnail` |
| `docs/<topic>`, `chore/<topic>` | documentation or maintenance |

Branch names describe the change, use lowercase words joined by hyphens, and stay under
eight words.

## Commits

- Use a prefix: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `ci`.
- Keep the whole message at eight words or fewer, prefix included:
  `fix: keep thumbnails when yt-dlp omits them`.
- Keep each commit at or under 100 inserted lines. Split big work by meaning
  (model, then API, then UI), never into "part 1" and "part 2".
- One clear change per commit. The message says exactly what changed.
- No generated trailers or tool signatures in commit messages.

## Pull requests

1. Branch from `develop` and open the pull request against `develop`.
2. Describe the change in plain technical terms: what changed, why, and how to test it.
3. Anything people can see needs screenshots of the changed area at desktop and phone
   width. Keep screenshots out of the repository; attach them to the pull request.
4. All checks pass: `gofmt -l .` prints nothing, `go vet ./...` and `go test ./...` pass.

## Code style

- Code is grouped by feature. Go packages live in `internal/<feature>`; dashboard code
  lives in `web/app/scripts/features/<feature>`.
- Standard library only on the Go side and plain JavaScript modules on the web side.
  Propose a new dependency in an issue before adding it.
- Name things for what they do, keep functions short, and comment the why, not the what.
- Shared helpers go in one place (`web/app/scripts/lib`, or the package that owns the
  concept). Do not copy code between features.
- Colors, sizes and radii come from `web/app/styles/tokens.css`. Do not hard-code them.

## Local setup

```sh
go run ./cmd/grabb -addr 127.0.0.1:7878
go test ./...
```

The dashboard is served from the embedded `web/app` folder, so restart the server to see
changes to HTML, CSS or JavaScript.

## Conduct

Everyone taking part follows the [Code of Conduct](CODE_OF_CONDUCT.md).
