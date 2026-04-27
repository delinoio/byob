# BYOB

BYOB means Bring Your Own Build tool.

This repository is a lightweight monorepo for building build-tool frameworks on top of the native TypeScript-Go compiler implementation. The current layout keeps integration boundaries explicit:

- `cmds/byob`: Go CLI entry point.
- `lint/`: Go runtime API for user-authored BYOB linter binaries.
- `lint/upstream`: Go helper for executable linters that run every pinned upstream rslint rule.
- `packages/byob`: TypeScript SDK surface for authoring build tools.
- `packages/rslint-bridge`: Go bridge for the pinned rslint all-rules engine.
- `packages/rslint-upstream`: Pinned rslint source snapshot.
- `packages/tsgo-bridge`: Go bridge that direct-links `github.com/microsoft/typescript-go`.
- `packages/typescript-go-upstream`: TypeScript-Go snapshot matched to rslint shims.
- `docs/`: Source of truth for project and domain contracts.

## Direct TypeScript-Go Link

`github.com/microsoft/typescript-go` currently keeps most compiler packages under `internal/`, so BYOB links through `packages/tsgo-bridge`.

The bridge module path is `github.com/microsoft/typescript-go/byobbridge`, which lets it import TypeScript-Go internals while exposing a small public API for the BYOB root module.

## Lint CLI

```sh
go run ./cmds/byob lint build --main ./path/to/main.go
go run ./cmds/byob lint run --main ./path/to/main.go -- --format json
```

`byob lint build` caches user linter binaries under the user cache directory. Passing `--out <dir>` exports the binary with `byob-lint-artifact.json` for sharing through CI artifacts or other team workflows.

The all-rules example builds a real linter that applies the pinned upstream rslint rule registry:

```sh
go run ./cmds/byob lint build --main examples/lint-all-rules/main.go
go run ./cmds/byob lint run --main examples/lint-all-rules/main.go -- --all-rules --type-check examples/lint-all-rules/fixtures
```

## Local Validation

```sh
go test ./...
go -C lint test ./...
go -C lint/upstream test ./...
go -C packages/rslint-bridge test ./...
go -C packages/tsgo-bridge test ./...
pnpm install
pnpm build
pnpm test
go run ./cmds/byob version
```
