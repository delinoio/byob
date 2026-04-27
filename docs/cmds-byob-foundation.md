# cmds-byob-foundation

## Scope
- Project/component: BYOB CLI.
- Canonical path: `cmds/byob`.

## Runtime and Language
- Runtime: Go CLI.
- Primary language: Go 1.26.

## Users and Operators
- Engineers validating BYOB workspace setup.
- Future build-tool authors using BYOB commands.

## Interfaces and Contracts
- Stable command identifier: `version`.
- `byob version` prints BYOB version, TypeScript-Go module path, and linked TypeScript-Go version.
- The CLI uses the bridge public API to prove direct TypeScript-Go linkage.

## Storage
- No persistent storage in the initial scaffold.

## Security
- The CLI must not execute arbitrary TypeScript-Go binaries for core compiler integration.
- Version output must avoid exposing local paths or environment values.

## Logging
- No structured runtime logging is required for the initial `version` command.
- Future commands should use structured logs for build lifecycle and compiler integration decisions.

## Build and Test
- Local validation: `go test ./...`.
- Runtime smoke check: `go run ./cmds/byob version`.

## Dependencies and Integrations
- Depends on `packages/tsgo-bridge` via `github.com/microsoft/typescript-go/byobbridge`.
- Does not import TypeScript-Go `internal/` packages directly.

## Change Triggers
- Update `docs/project-byob.md` and this file when command identifiers or CLI output contracts change.
- Update `docs/packages-tsgo-bridge-foundation.md` when CLI changes require new bridge surface area.

## References
- `docs/project-byob.md`
- `docs/packages-tsgo-bridge-foundation.md`
- `docs/domain-template.md`

