# Project: byob

## Goal
Build a framework for authoring, building, and sharing custom build tools on top of TypeScript-Go while keeping integration boundaries explicit and replaceable.

## Project ID
`byob`

## Domain Ownership Map
- `cmds/byob`: BYOB command-line entry point.
- `lint`: BYOB Go lint runtime API for user-authored linter binaries.
- `lint/upstream`: BYOB helper module for executable all-rules rslint linters.
- `packages/byob`: TypeScript SDK for build-tool definitions.
- `packages/rslint-bridge`: Go bridge that exposes pinned rslint all-rules execution through BYOB-owned public types.
- `packages/rslint-upstream`: Pinned upstream rslint Go source snapshot.
- `packages/tsgo-bridge`: Go bridge that direct-links TypeScript-Go internals through an allowed module path.
- `packages/typescript-go-upstream`: Pinned TypeScript-Go source snapshot used by rslint generated shims.
- `docs`: Repository and domain contracts.
- `scripts`: Local development helpers.

## Domain Contract Documents
- `docs/cmds-byob-foundation.md`
- `docs/lint-byob-runtime-foundation.md`
- `docs/packages-byob-sdk-foundation.md`
- `docs/packages-rslint-bridge-foundation.md`
- `docs/packages-rslint-upstream-foundation.md`
- `docs/packages-tsgo-bridge-foundation.md`
- `docs/packages-typescript-go-upstream-foundation.md`

## Cross-Domain Invariants
- Root BYOB Go code imports `github.com/microsoft/typescript-go/byobbridge`, not `github.com/microsoft/typescript-go/internal/...`.
- The bridge module path remains `github.com/microsoft/typescript-go/byobbridge` so Go `internal/` import rules allow direct TypeScript-Go linkage.
- The TypeScript SDK package name remains `@delinoio/byob`.
- Stable CLI command identifiers include `version` and `lint`.
- Stable BYOB tool-kind identifiers include `lint`, `fmt`, and `transform`; only `lint` is implemented initially.
- BYOB root code must not import `github.com/web-infra-dev/rslint/internal/...` or shell out to `rslint` for core lint integration.
- Only `packages/rslint-bridge` may import pinned rslint internals, and it must expose BYOB-owned public types.

## Change Policy
- Update this index and the relevant domain contract whenever CLI commands, SDK types, bridge imports, or repository structure change.
- Update root `AGENTS.md` when repository-wide validation or TypeScript-Go integration rules change.

## References
- `docs/README.md`
- `docs/project-template.md`
- `docs/domain-template.md`
