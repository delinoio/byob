# BYOB

BYOB means Bring Your Own Build tool.

This repository is a lightweight monorepo for building build-tool frameworks on top of the native TypeScript-Go compiler implementation. The initial scaffold keeps the runtime small:

- `cmds/byob`: Go CLI entry point.
- `packages/byob`: TypeScript SDK surface for authoring build tools.
- `packages/tsgo-bridge`: Go bridge that direct-links `github.com/microsoft/typescript-go`.
- `docs/`: Source of truth for project and domain contracts.

## Direct TypeScript-Go Link

`github.com/microsoft/typescript-go` currently keeps most compiler packages under `internal/`, so BYOB links through `packages/tsgo-bridge`.

The bridge module path is `github.com/microsoft/typescript-go/byobbridge`, which lets it import TypeScript-Go internals while exposing a small public API for the BYOB root module.

## Local Validation

```sh
go test ./...
go -C packages/tsgo-bridge test ./...
pnpm install
pnpm build
pnpm test
go run ./cmds/byob version
```

