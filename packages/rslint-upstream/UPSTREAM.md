# rslint upstream snapshot

- Repository: `https://github.com/web-infra-dev/rslint`
- Pinned commit: `b986707ef58537329229645b77d35d93e3234ca6`
- Snapshot contents: Go `internal/` packages and generated `shim/` modules needed by BYOB's rslint bridge.

## Sync

```sh
tmp="$(mktemp -d)"
git clone --depth 1 https://github.com/web-infra-dev/rslint.git "$tmp/rslint"
git -C "$tmp/rslint" checkout b986707ef58537329229645b77d35d93e3234ca6
rsync -a --delete "$tmp/rslint/internal/" packages/rslint-upstream/internal/
rsync -a --delete "$tmp/rslint/shim/" packages/rslint-upstream/shim/
cp "$tmp/rslint/go.mod" packages/rslint-upstream/go.mod
cp "$tmp/rslint/go.sum" packages/rslint-upstream/go.sum
```
