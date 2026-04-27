# packages-rslint-upstream-foundation

## Scope
- Project/component: pinned rslint upstream source snapshot.
- Canonical path: `packages/rslint-upstream`.

## Runtime and Language
- Runtime: Go module.
- Primary language: Go 1.26.

## Users and Operators
- BYOB rslint bridge maintainers.
- Engineers auditing the exact upstream rule implementation used by BYOB lint examples.

## Interfaces and Contracts
- Module path: `github.com/web-infra-dev/rslint`.
- Pinned commit: `b986707ef58537329229645b77d35d93e3234ca6`.
- Snapshot contents include upstream Go `internal/` packages and generated `shim/` modules.
- `UPSTREAM.md` records the source repo, commit, and sync command.

## Storage
- The snapshot owns no runtime storage.

## Security
- Root BYOB code must not import this module's `internal/` packages directly.
- BYOB access to this snapshot must flow through `packages/rslint-bridge`.

## Logging
- No BYOB-owned logging behavior lives in the snapshot.

## Build and Test
- The snapshot is validated through `go -C packages/rslint-bridge test ./...`.
- Full upstream test execution is outside the BYOB v1 validation contract.

## Dependencies and Integrations
- Used by `packages/rslint-bridge`.
- Uses local shim module replacements under `packages/rslint-upstream/shim`.
- Requires the matching `packages/typescript-go-upstream` snapshot when built from the BYOB workspace.

## Change Triggers
- Update this file and `UPSTREAM.md` when changing snapshot contents or upstream commit.
- Update bridge tests when upstream rule count or representative rule names change.

## References
- `docs/project-byob.md`
- `docs/packages-rslint-bridge-foundation.md`
