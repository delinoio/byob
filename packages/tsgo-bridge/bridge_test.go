package byobbridge

import "testing"

func TestInfoReportsTypeScriptGoLink(t *testing.T) {
	info := Info()
	if info.Module != TypeScriptGoModule {
		t.Fatalf("unexpected module: got=%s want=%s", info.Module, TypeScriptGoModule)
	}
	if info.Version == "" {
		t.Fatal("expected linked TypeScript-Go version")
	}
}
