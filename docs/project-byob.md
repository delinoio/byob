# Project: byob

## Goal
Build a framework for authoring custom build tools on top of TypeScript-Go while keeping integration boundaries explicit and replaceable.

## Project ID
`byob`

## Domain Ownership Map
- `cmds/byob`: BYOB command-line entry point.
- `packages/byob`: TypeScript SDK for build-tool definitions.
- `packages/tsgo-bridge`: Go bridge that direct-links TypeScript-Go internals through an allowed module path.
- `docs`: Repository and domain contracts.
- `scripts`: Local development helpers.

## Domain Contract Documents
- `docs/cmds-byob-foundation.md`
- `docs/packages-byob-sdk-foundation.md`
- `docs/packages-tsgo-bridge-foundation.md`

## Cross-Domain Invariants
- Root BYOB Go code imports `github.com/microsoft/typescript-go/byobbridge`, not `github.com/microsoft/typescript-go/internal/...`.
- The bridge module path remains `github.com/microsoft/typescript-go/byobbridge` so Go `internal/` import rules allow direct TypeScript-Go linkage.
- The TypeScript SDK package name remains `@delinoio/byob`.
- The initial CLI command identifier is `version`.

## Change Policy
- Update this index and the relevant domain contract whenever CLI commands, SDK types, bridge imports, or repository structure change.
- Update root `AGENTS.md` when repository-wide validation or TypeScript-Go integration rules change.

## References
- `docs/README.md`
- `docs/project-template.md`
- `docs/domain-template.md`

