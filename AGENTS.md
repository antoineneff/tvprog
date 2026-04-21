# AGENTS.md

## Build & Run

- Build: `go build -o tvprog main.go`
- Run: `./tvprog` (serves on `:3000`)
- No tests exist in this repo.

## Environment

- Go version is pinned via `mise.toml` to **1.25.6**. Prefer using `mise` or matching this version locally.
- `go.mod` module name is `tvprog`; internal imports use `tvprog/pkg/...`.

## Runtime Behavior

- The app fetches French TV program data from `https://xmltvfr.fr/xmltv/xmltv_tnt.zip` at startup and every day at 04:00.
- Two endpoints:
  - `GET /` — plain text formatted table
  - `GET /json` — JSON representation of the day's programs

## Deployment

- `Dockerfile` produces a minimal Alpine image. Build artifact is `/bin/api`.
- `GIN_MODE=release` is set in the Docker image.

## Notes

- `.gitignore` excludes the `tvprog` binary.
- No CI, Makefile, or test suite exists.
