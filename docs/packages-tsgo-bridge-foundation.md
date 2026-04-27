# packages-tsgo-bridge-foundation

## Scope
- Project/component: TypeScript-Go bridge.
- Canonical path: `packages/tsgo-bridge`.

## Runtime and Language
- Runtime: Go module.
- Primary language: Go 1.26.

## Users and Operators
- BYOB Go components that need direct TypeScript-Go linkage.
- Maintainers updating pinned TypeScript-Go integration points.

## Interfaces and Contracts
- Module path: `github.com/microsoft/typescript-go/byobbridge`.
- Pinned dependency remains declared as `github.com/microsoft/typescript-go v0.0.0-20260424234512-515d036f927a`; the BYOB workspace replaces it with `packages/typescript-go-upstream` to keep the rslint shim and direct bridge on one internal API surface.
- Public API includes `Info()` and `LinkedVersion()`.
- The bridge is the only BYOB-owned module allowed to import `github.com/microsoft/typescript-go/internal/...`.

## Storage
- No persistent storage.

## Security
- The bridge links Go packages directly and must not invoke external `tsgo` processes for core compiler access.
- Public bridge APIs should expose stable BYOB-owned shapes rather than leaking TypeScript-Go internal types unless intentionally documented.

## Logging
- No logging is required for the initial version-link smoke API.
- Future compiler integration APIs should accept logger or context hooks rather than writing directly to global output.

## Build and Test
- Local validation: `go -C packages/tsgo-bridge test ./...`.
- Root integration validation: `go test ./...`.

## Dependencies and Integrations
- Directly imports TypeScript-Go internals from `github.com/microsoft/typescript-go/internal/core`.
- Uses `packages/typescript-go-upstream` through local replacement in this repository.
- Consumed by `cmds/byob` through public bridge APIs.

## Change Triggers
- Update `docs/project-byob.md` and this file when the bridge module path, pinned TypeScript-Go version, or public bridge API changes.
- Update `docs/cmds-byob-foundation.md` when bridge changes affect CLI output or command behavior.
- Update `docs/packages-typescript-go-upstream-foundation.md` when the workspace TypeScript-Go snapshot changes.

## References
- `docs/project-byob.md`
- `docs/cmds-byob-foundation.md`
- `docs/domain-template.md`
