# Security Policy

## Supported versions

Security fixes go into the latest release on `main`.

## Reporting a vulnerability

Please do not open a public issue. Use
[private vulnerability reporting](https://github.com/isonlycoolie/grabb/security/advisories/new)
on GitHub and include the steps to reproduce, the affected version and the impact you see.

You can expect an acknowledgement within three days and a plan for a fix within two weeks.

## How grabb protects you

- The server listens on `127.0.0.1` unless you choose another address with `-addr`.
- Write requests from other websites are refused, and the server only answers when it is
  addressed by IP or `localhost`, which blocks DNS rebinding.
- Links must be absolute `http` or `https` URLs and are passed to yt-dlp after `--`, so a
  link can never be read as a command-line option.
- The dashboard sends a strict Content Security Policy and loads no third-party code.

If you expose grabb on your network with `-addr 0.0.0.0:7878`, anyone on that network can
use it. Do this only on networks you trust.
