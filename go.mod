module github.com/delinoio/byob

go 1.26

require (
	github.com/delinoio/byob/lint v0.0.0
	github.com/microsoft/typescript-go/byobbridge v0.0.0
)

require (
	github.com/go-json-experiment/json v0.0.0-20260214004413-d219187c3433 // indirect
	github.com/microsoft/typescript-go v0.0.0-20260424234512-515d036f927a // indirect
	golang.org/x/sync v0.20.0 // indirect
)

replace github.com/microsoft/typescript-go/byobbridge => ./packages/tsgo-bridge

replace github.com/delinoio/byob/lint => ./lint
