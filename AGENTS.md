### Instructions

- Use the `docs/` directory as the source of truth for project contracts and implementation documents.
- Keep repository-wide rules in this file and domain-specific contract rules in `docs/`.
- Write all source code and comments in English.
- Prefer typed constants and explicit interfaces over free-form strings when contracts may become stable.
- Go code in this repository targets Go 1.26.
- BYOB must direct-link `github.com/microsoft/typescript-go` through `packages/tsgo-bridge`; do not shell out to `tsgo` for core compiler integration.
- Root BYOB code must import the bridge public API instead of importing `github.com/microsoft/typescript-go/internal/...` directly.
- If command, SDK, bridge, or repository structure changes, update the relevant `docs/project-*.md` and `docs/*-foundation.md` files in the same change.

### Monorepo Structure Map

- `docs/`: Source of truth for project contracts and repository documentation.
- `cmds/`: Go command tools.
- `packages/`: TypeScript packages and Go bridge modules.
- `scripts/`: Local development helper scripts.

### Validation

- Run `go test ./...` when root Go code changes.
- Run `go -C packages/tsgo-bridge test ./...` when the TypeScript-Go bridge changes.
- Run `pnpm build` and `pnpm test` when TypeScript packages change.

