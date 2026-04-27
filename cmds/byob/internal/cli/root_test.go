package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteRequiresCommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := execute([]string{}, stdout, stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got=%d", code)
	}
	if !strings.Contains(stderr.String(), "usage:") {
		t.Fatalf("expected usage output, got=%s", stderr.String())
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := execute([]string{"unknown"}, stdout, stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got=%d", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("expected unknown command error, got=%s", stderr.String())
	}
}

func TestVersionReportsDirectLink(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := execute([]string{"version"}, stdout, stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got=%d stderr=%s", code, stderr.String())
	}

	output := stdout.String()
	for _, expected := range []string{
		"byob version:",
		"typescript-go module: github.com/microsoft/typescript-go",
		"typescript-go version:",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got=%s", expected, output)
		}
	}
}
