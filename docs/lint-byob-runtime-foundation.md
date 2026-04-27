# lint-byob-runtime-foundation

## Scope
- Project/component: BYOB Go lint runtime API.
- Canonical path: `lint`.
- Module path: `github.com/delinoio/byob/lint`.

## Runtime and Language
- Runtime: Go module.
- Primary language: Go 1.26.

## Users and Operators
- User linter authors importing BYOB runtime types.
- BYOB CLI maintainers building and executing user linter binaries.

## Interfaces and Contracts
- Public compatibility constants include `RuntimeVersion` and `RslintCompatVersion`.
- `RslintCompatVersion` tracks the upstream rslint rule-framework reference version; the initial reference is `0.5.0`.
- Public typed tool-kind constants include `lint`, `fmt`, and `transform`; only `lint` has CLI support initially.
- Public rule APIs include diagnostic severities, text ranges, nodes, source files, programs, rule listeners, rule definitions, diagnostics, fixes, suggestions, and report callbacks.
- `CreateRule` and `CreateTypeScriptRule` prefix TypeScript-ESLint rules with `@typescript-eslint/`, matching the rslint v0.5.0 convention.
- `github.com/delinoio/byob/lint/upstream` provides `RunAllRulesCLI(args)` and `RunAllRulesCLIWithIO(...)` for executable linter mains that run every pinned upstream rslint rule.

## Storage
- The runtime package owns no persistent storage.
- Built user linter binaries are owned by the CLI cache contract in `docs/cmds-byob-foundation.md`.

## Security
- The runtime package must not import `github.com/web-infra-dev/rslint/internal/...`.
- The runtime package must not shell out to `rslint` or `tsgo`.
- Future TypeScript-Go access must flow through BYOB-owned bridge/public APIs instead of direct root imports of TypeScript-Go internals.
- Upstream all-rules execution must flow through `packages/rslint-bridge`; user linter mains must not import rslint internals.

## Logging
- The runtime package does not log directly.
- Future executable lint engines should accept explicit logger or context hooks.

## Build and Test
- Local validation: `go -C lint test ./...`.
- Upstream helper validation: `go -C lint/upstream test ./...`.
- Root CLI integration validation: `go test ./...`.

## Dependencies and Integrations
- The initial runtime API has no external module dependencies.
- The root CLI imports the runtime for compatibility metadata when writing lint artifacts.
- The upstream helper module depends on `packages/rslint-bridge` and is intentionally separate from the lightweight root `lint` module.
- Future rslint rule ports should be source-level BYOB-owned ports or syncs, not imports from rslint internal packages.

## Change Triggers
- Update `docs/project-byob.md` and this file when the runtime module path, public rule API, compatibility version, or tool-kind constants change.
- Update `docs/cmds-byob-foundation.md` when runtime changes affect CLI build, run, cache, or artifact behavior.
- Update `docs/packages-tsgo-bridge-foundation.md` when runtime changes require new bridge surface area.
- Update `docs/packages-rslint-bridge-foundation.md` when upstream helper behavior changes.

## References
- `docs/project-byob.md`
- `docs/cmds-byob-foundation.md`
- `docs/packages-rslint-bridge-foundation.md`
- `docs/domain-template.md`
