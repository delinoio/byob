# cmds-byob-foundation

## Scope
- Project/component: BYOB CLI.
- Canonical path: `cmds/byob`.

## Runtime and Language
- Runtime: Go CLI.
- Primary language: Go 1.26.

## Users and Operators
- Engineers validating BYOB workspace setup.
- Build-tool authors building and sharing user-authored BYOB tools.
- Teams running cached BYOB tool binaries in local development or CI.

## Interfaces and Contracts
- Stable command identifiers: `version`, `lint`.
- `byob version` prints BYOB version, TypeScript-Go module path, and linked TypeScript-Go version.
- The CLI uses the bridge public API to prove direct TypeScript-Go linkage.
- `byob lint build --main <path> [--out <dir>] [--target <goos/goarch>] [--force]` builds a user linter package from an explicit Go `main.go` file.
- `byob lint run --main <path> [--target <goos/goarch>] [--] <args...>` builds or reuses the cached linter binary, forwards stdio, and returns the linter process exit code.
- `lint run` defaults to the host `GOOS/GOARCH` and rejects non-host targets because cross-built binaries cannot be executed locally.
- `--out` exports the binary and `byob-lint-artifact.json`; external uploads are outside the v1 CLI contract.

## Storage
- Lint binaries are cached under `os.UserCacheDir()/byob/lint/<cache-key>/`.
- The lint cache key includes target, BYOB version, linked TypeScript-Go version, rslint compatibility version, resolved main package, local source inputs, and Go module/workspace dependency files.

## Security
- The CLI must not execute arbitrary TypeScript-Go binaries for core compiler integration.
- Version output must avoid exposing local paths or environment values.
- `lint run` executes only the user linter binary built from the explicit `--main` package.
- BYOB must not shell out to `rslint` or `tsgo` for core lint/compiler integration.

## Logging
- No structured runtime logging is required for the initial `version` command.
- `lint build` and cache lifecycle messages are written to stderr so linter stdout remains available to callers.
- Future commands should use structured logs for build lifecycle and compiler integration decisions.

## Build and Test
- Local validation: `go test ./...`.
- Lint runtime validation: `go -C lint test ./...`.
- Runtime smoke check: `go run ./cmds/byob version`.
- Lint smoke check: `go run ./cmds/byob lint build --main <linter-main.go>`.

## Dependencies and Integrations
- Depends on `packages/tsgo-bridge` via `github.com/microsoft/typescript-go/byobbridge`.
- Depends on `lint` via `github.com/delinoio/byob/lint` for runtime compatibility metadata.
- Does not import TypeScript-Go `internal/` packages directly.
- Does not import `github.com/web-infra-dev/rslint/internal/...`.

## Change Triggers
- Update `docs/project-byob.md` and this file when command identifiers or CLI output contracts change.
- Update `docs/lint-byob-runtime-foundation.md` when lint runtime public API assumptions change.
- Update `docs/packages-tsgo-bridge-foundation.md` when CLI changes require new bridge surface area.

## References
- `docs/project-byob.md`
- `docs/lint-byob-runtime-foundation.md`
- `docs/packages-tsgo-bridge-foundation.md`
- `docs/domain-template.md`
