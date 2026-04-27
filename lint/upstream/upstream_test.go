package upstream

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAllRulesCLIReportsDiagnostics(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "example.ts")
	if err := os.WriteFile(sourcePath, []byte(`console.log("hello")`), 0644); err != nil {
		t.Fatal(err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := RunAllRulesCLIWithIO([]string{"--format", "jsonline", sourcePath}, nil, stdout, stderr)
	if code == 0 {
		t.Fatalf("expected lint diagnostics to produce nonzero exit, stderr=%s stdout=%s", stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), `"ruleName":"no-console"`) {
		t.Fatalf("expected no-console diagnostic, stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
}
