# packages-rslint-bridge-foundation

## Scope
- Project/component: BYOB rslint bridge.
- Canonical path: `packages/rslint-bridge`.

## Runtime and Language
- Runtime: Go module.
- Primary language: Go 1.26.

## Users and Operators
- BYOB lint helper packages that need the pinned upstream rslint engine.
- Maintainers syncing BYOB's rslint integration to a new upstream commit.

## Interfaces and Contracts
- Module path: `github.com/web-infra-dev/rslint/byobbridge`.
- Public API includes `Info()`, `UpstreamCommit()`, `AllRuleNames()`, `RunAllRules(options)`, `RunOptions`, and typed `OutputFormat*` constants.
- Pinned rslint commit: `b986707ef58537329229645b77d35d93e3234ca6`.
- `AllRuleNames` returns every rule registered by upstream rslint `RegisterAllRules`.
- `RunAllRules` enables all registered rules at error severity by default and accepts explicit files, config path, format, fix, type-check, quiet, max-warnings, rule override, and stdio options.

## Storage
- The bridge owns no persistent storage.
- `RunAllRules` writes only when `RunOptions.Fix` is true.

## Security
- The bridge is the only BYOB-owned package allowed to import `github.com/web-infra-dev/rslint/internal/...`.
- The bridge must not shell out to `rslint` or `tsgo`.
- Public APIs expose BYOB-owned shapes rather than rslint internal types.

## Logging
- Runtime diagnostics and lifecycle errors go to the provided stdio writers.
- No global logger is used.

## Build and Test
- Local validation: `go -C packages/rslint-bridge test ./...`.
- Integration validation: `go test ./...`.

## Dependencies and Integrations
- Depends on `packages/rslint-upstream` for the pinned upstream rslint snapshot.
- Depends on `packages/typescript-go-upstream` for the TypeScript-Go internals matched to rslint's generated shims.
- Consumed by `github.com/delinoio/byob/lint/upstream`.

## Change Triggers
- Update `docs/project-byob.md`, this file, and `packages/rslint-upstream/UPSTREAM.md` when the pinned rslint commit changes.
- Update `docs/packages-typescript-go-upstream-foundation.md` when the matching TypeScript-Go snapshot changes.

## References
- `docs/project-byob.md`
- `docs/lint-byob-runtime-foundation.md`
- `docs/packages-rslint-upstream-foundation.md`
- `docs/packages-typescript-go-upstream-foundation.md`
