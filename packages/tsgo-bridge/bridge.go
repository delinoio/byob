package byobbridge

import "github.com/microsoft/typescript-go/internal/core"

const TypeScriptGoModule = "github.com/microsoft/typescript-go"

type LinkInfo struct {
	Module  string
	Version string
}

func Info() LinkInfo {
	return LinkInfo{
		Module:  TypeScriptGoModule,
		Version: LinkedVersion(),
	}
}

func LinkedVersion() string {
	return core.Version()
}
