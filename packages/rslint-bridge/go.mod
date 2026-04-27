module github.com/web-infra-dev/rslint/byobbridge

go 1.26.0

require (
	github.com/microsoft/typescript-go/shim/ast v0.0.0
	github.com/microsoft/typescript-go/shim/compiler v0.0.0
	github.com/microsoft/typescript-go/shim/core v0.0.0
	github.com/microsoft/typescript-go/shim/scanner v0.0.0
	github.com/microsoft/typescript-go/shim/tspath v0.0.0
	github.com/microsoft/typescript-go/shim/vfs v0.0.0
	github.com/microsoft/typescript-go/shim/vfs/osvfs v0.0.0
	github.com/web-infra-dev/rslint v0.0.0
)

require (
	github.com/bmatcuk/doublestar/v4 v4.10.0 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/go-json-experiment/json v0.0.0-20260214004413-d219187c3433 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/microsoft/typescript-go v0.0.0-20260313230633-c0e5d35a6f8f // indirect
	github.com/microsoft/typescript-go/shim/bundled v0.0.0 // indirect
	github.com/microsoft/typescript-go/shim/checker v0.0.0 // indirect
	github.com/microsoft/typescript-go/shim/evaluator v0.0.0-00010101000000-000000000000 // indirect
	github.com/microsoft/typescript-go/shim/tsoptions v0.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/tailscale/hujson v0.0.0-20250605163823-992244df8c5a // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
)

replace github.com/web-infra-dev/rslint => ../rslint-upstream

replace github.com/microsoft/typescript-go => ../typescript-go-upstream

replace github.com/microsoft/typescript-go/shim/api => ../rslint-upstream/shim/api

replace github.com/microsoft/typescript-go/shim/api/encoder => ../rslint-upstream/shim/api/encoder

replace github.com/microsoft/typescript-go/shim/ast => ../rslint-upstream/shim/ast

replace github.com/microsoft/typescript-go/shim/bundled => ../rslint-upstream/shim/bundled

replace github.com/microsoft/typescript-go/shim/checker => ../rslint-upstream/shim/checker

replace github.com/microsoft/typescript-go/shim/collections => ../rslint-upstream/shim/collections

replace github.com/microsoft/typescript-go/shim/compiler => ../rslint-upstream/shim/compiler

replace github.com/microsoft/typescript-go/shim/core => ../rslint-upstream/shim/core

replace github.com/microsoft/typescript-go/shim/evaluator => ../rslint-upstream/shim/evaluator

replace github.com/microsoft/typescript-go/shim/jsonrpc => ../rslint-upstream/shim/jsonrpc

replace github.com/microsoft/typescript-go/shim/ls => ../rslint-upstream/shim/ls

replace github.com/microsoft/typescript-go/shim/lsp/lsproto => ../rslint-upstream/shim/lsp/lsproto

replace github.com/microsoft/typescript-go/shim/project => ../rslint-upstream/shim/project

replace github.com/microsoft/typescript-go/shim/scanner => ../rslint-upstream/shim/scanner

replace github.com/microsoft/typescript-go/shim/tsoptions => ../rslint-upstream/shim/tsoptions

replace github.com/microsoft/typescript-go/shim/tspath => ../rslint-upstream/shim/tspath

replace github.com/microsoft/typescript-go/shim/vfs => ../rslint-upstream/shim/vfs

replace github.com/microsoft/typescript-go/shim/vfs/cachedvfs => ../rslint-upstream/shim/vfs/cachedvfs

replace github.com/microsoft/typescript-go/shim/vfs/osvfs => ../rslint-upstream/shim/vfs/osvfs
