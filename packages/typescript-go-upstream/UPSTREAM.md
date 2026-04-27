# TypeScript-Go upstream snapshot

- Repository: `https://github.com/microsoft/typescript-go`
- Pinned commit: `93ae10c465dce1c1ac5c0d205c4b19597fc969f6`
- Purpose: matching TypeScript-Go internals for the pinned rslint generated shim modules.
- Snapshot contents: root `go.mod`, `go.sum`, and the `internal/` packages required by rslint shims, excluding TypeScript-Go test-only internals.

## Sync

```sh
tmp="$(mktemp -d)"
git clone --depth 1 https://github.com/web-infra-dev/rslint.git "$tmp/rslint"
git -C "$tmp/rslint" submodule update --init typescript-go
rsync -a --delete \
  --exclude 'internal/fourslash' \
  --exclude 'internal/testutil' \
  --exclude 'internal/testrunner' \
  "$tmp/rslint/typescript-go/internal/" packages/typescript-go-upstream/internal/
cp "$tmp/rslint/typescript-go/go.mod" packages/typescript-go-upstream/go.mod
cp "$tmp/rslint/typescript-go/go.sum" packages/typescript-go-upstream/go.sum
```
