# Repository Guidelines — isucrud

ISUCON-oriented CRUD analyzer / static-analysis CLI for Go. Releases are
packaged with goreleaser.

> Agent configuration is managed via [apm](https://github.com/microsoft/apm).
> Common conventions live in `mazrean/apm-plackage/common`; Go and goreleaser
> rules come from `mazrean/apm-plackage/{go,goreleaser}`. Run `apm install`
> to materialise locally.

## Build & Test

- `go test -v ./...`
- `go build ./...`
- `golangci-lint run`
- `goreleaser release --snapshot --clean` — cross-build snapshot for local testing

## Conventions

- Specs go under `specs/`; use `mazrean/agent-skills/skills/writing-*`.
- Commit using Conventional Commits (`committing-code` skill).
- Use the Go 1.24+ `tool` directive for build tools; see `using-go-tool-directive` skill.
- Releases via goreleaser; tag with `git tag vX.Y.Z` and let CI publish.
