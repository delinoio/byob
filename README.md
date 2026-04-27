# BYOB

BYOB means Bring Your Own Build tool.

This repository is a lightweight monorepo for building build-tool frameworks on top of the native TypeScript-Go compiler implementation. The current layout keeps integration boundaries explicit:

- `cmds/byob`: Go CLI entry point.
- `lint/`: Go runtime API for user-authored BYOB linter binaries.
- `packages/byob`: TypeScript SDK surface for authoring build tools.
- `packages/tsgo-bridge`: Go bridge that direct-links `github.com/microsoft/typescript-go`.
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

## Local Validation

```sh
go test ./...
go -C lint test ./...
go -C packages/tsgo-bridge test ./...
pnpm install
pnpm build
pnpm test
go run ./cmds/byob version
```
