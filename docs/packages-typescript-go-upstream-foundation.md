# packages-typescript-go-upstream-foundation

## Scope
- Project/component: pinned TypeScript-Go source snapshot for rslint shims.
- Canonical path: `packages/typescript-go-upstream`.

## Runtime and Language
- Runtime: Go module.
- Primary language: Go 1.26.

## Users and Operators
- Maintainers of the BYOB rslint integration.
- Engineers validating TypeScript-Go shim compatibility.

## Interfaces and Contracts
- Module path: `github.com/microsoft/typescript-go`.
- Pinned commit: `93ae10c465dce1c1ac5c0d205c4b19597fc969f6`.
- Snapshot contents include root `go.mod`, `go.sum`, and TypeScript-Go `internal/` packages required by rslint shims.
- `UPSTREAM.md` records the source repo, commit, and sync command.

## Storage
- The snapshot owns no runtime storage.

## Security
- Root BYOB code must continue importing TypeScript-Go through `packages/tsgo-bridge`.
- This snapshot exists to satisfy vendored upstream rslint shim imports and should not become a new root BYOB integration surface.

## Logging
- No BYOB-owned logging behavior lives in the snapshot.

## Build and Test
- Validated through `go -C packages/rslint-bridge test ./...` and `go -C packages/tsgo-bridge test ./...`.

## Dependencies and Integrations
- Used by rslint generated shim modules and the BYOB TypeScript-Go bridge through local replaces.
- Must remain aligned with `packages/rslint-upstream`.

## Change Triggers
- Update this file and `UPSTREAM.md` when changing snapshot contents or pinned commit.
- Update `docs/packages-tsgo-bridge-foundation.md` when the TypeScript-Go bridge dependency pin changes.

## References
- `docs/project-byob.md`
- `docs/packages-rslint-bridge-foundation.md`
- `docs/packages-rslint-upstream-foundation.md`
