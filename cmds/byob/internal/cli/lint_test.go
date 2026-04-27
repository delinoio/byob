package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	byoblint "github.com/delinoio/byob/lint"
)

func TestLintRequiresSubcommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := execute([]string{"lint"}, stdout, stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got=%d", code)
	}
	if !strings.Contains(stderr.String(), "usage: byob lint <command>") {
		t.Fatalf("expected lint usage, got=%s", stderr.String())
	}
}

func TestLintUnknownSubcommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := execute([]string{"lint", "unknown"}, stdout, stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got=%d", code)
	}
	if !strings.Contains(stderr.String(), "unknown lint command") {
		t.Fatalf("expected unknown lint command error, got=%s", stderr.String())
	}
}

func TestLintBuildRequiresMain(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := execute([]string{"lint", "build"}, stdout, stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got=%d", code)
	}
	if !strings.Contains(stderr.String(), "requires --main") {
		t.Fatalf("expected missing main error, got=%s", stderr.String())
	}
}

func TestLintRunRejectsNonHostTarget(t *testing.T) {
	ctx := testCommandContext(t)
	code, _, stderr := executeTestCommand(t, ctx, []string{
		"lint",
		"run",
		"--main",
		"does-not-need-to-exist.go",
		"--target",
		nonHostTarget(),
	})
	if code != 2 {
		t.Fatalf("expected exit code 2, got=%d", code)
	}
	if !strings.Contains(stderr, "cannot execute on host") {
		t.Fatalf("expected host target error, got=%s", stderr)
	}
}

func TestLintBuildCachesAndExportsArtifact(t *testing.T) {
	mainPath := writeTestLinter(t)
	ctx := testCommandContext(t)

	code, _, stderr := executeTestCommand(t, ctx, []string{"lint", "build", "--main", mainPath})
	if code != 0 {
		t.Fatalf("expected first build to succeed, code=%d stderr=%s", code, stderr)
	}
	if !strings.Contains(stderr, "built lint binary:") {
		t.Fatalf("expected build message, got=%s", stderr)
	}

	code, _, stderr = executeTestCommand(t, ctx, []string{"lint", "build", "--main", mainPath})
	if code != 0 {
		t.Fatalf("expected cached build to succeed, code=%d stderr=%s", code, stderr)
	}
	if !strings.Contains(stderr, "using cached lint binary:") {
		t.Fatalf("expected cache message, got=%s", stderr)
	}

	code, _, stderr = executeTestCommand(t, ctx, []string{"lint", "build", "--main", mainPath, "--force"})
	if code != 0 {
		t.Fatalf("expected forced build to succeed, code=%d stderr=%s", code, stderr)
	}
	if !strings.Contains(stderr, "built lint binary:") {
		t.Fatalf("expected forced build message, got=%s", stderr)
	}

	outDir := t.TempDir()
	code, _, stderr = executeTestCommand(t, ctx, []string{"lint", "build", "--main", mainPath, "--out", outDir})
	if code != 0 {
		t.Fatalf("expected export build to succeed, code=%d stderr=%s", code, stderr)
	}

	manifestPath := filepath.Join(outDir, "byob-lint-artifact.json")
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("expected manifest: %v", err)
	}
	var manifest lintArtifactManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("invalid manifest: %v", err)
	}
	if manifest.Tool != string(byoblint.ToolKindLint) {
		t.Fatalf("unexpected manifest tool: %s", manifest.Tool)
	}
	if manifest.Target != hostTarget().String() {
		t.Fatalf("unexpected manifest target: %s", manifest.Target)
	}
	if manifest.SourceMain != mainPath {
		t.Fatalf("unexpected manifest source main: %s", manifest.SourceMain)
	}
	if manifest.BYOBVersion != byobVersion {
		t.Fatalf("unexpected manifest BYOB version: %s", manifest.BYOBVersion)
	}
	if manifest.RslintCompatVersion == "" {
		t.Fatal("expected manifest rslint compatibility version")
	}
	if _, err := os.Stat(filepath.Join(outDir, manifest.Binary)); err != nil {
		t.Fatalf("expected exported binary: %v", err)
	}
}

func TestLintRunBuildsForwardsArgsAndExitCode(t *testing.T) {
	mainPath := writeTestLinter(t)
	ctx := testCommandContext(t)
	ctx.env = withEnv(ctx.env, "BYOB_TEST_EXIT_7", "1")

	code, stdout, stderr := executeTestCommand(t, ctx, []string{
		"lint",
		"run",
		"--main",
		mainPath,
		"--",
		"--alpha",
		"value",
	})
	if code != 7 {
		t.Fatalf("expected linter exit code 7, got=%d stderr=%s", code, stderr)
	}
	if !strings.Contains(stdout, "args=--alpha,value") {
		t.Fatalf("expected forwarded args in stdout, got=%s", stdout)
	}
}

func TestLintRunExampleUpstreamAllRules(t *testing.T) {
	root := repoRoot(t)
	mainPath := filepath.Join(root, "examples", "lint-all-rules", "main.go")
	fixturePath := filepath.Join(root, "examples", "lint-all-rules", "fixtures", "example.ts")
	ctx := testCommandContext(t)

	code, stdout, stderr := executeTestCommand(t, ctx, []string{
		"lint",
		"run",
		"--main",
		mainPath,
		"--",
		"--all-rules",
		"--format",
		"jsonline",
		fixturePath,
	})
	if code == 0 {
		t.Fatalf("expected upstream diagnostics to produce nonzero exit, stderr=%s stdout=%s", stderr, stdout)
	}
	if !strings.Contains(stdout, `"ruleName":"no-console"`) {
		t.Fatalf("expected no-console diagnostic, stderr=%s stdout=%s", stderr, stdout)
	}
}

func executeTestCommand(t *testing.T, ctx commandContext, args []string) (int, string, string) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	ctx.stdout = stdout
	ctx.stderr = stderr
	code := executeWithContext(args, ctx)
	return code, stdout.String(), stderr.String()
}

func testCommandContext(t *testing.T) commandContext {
	t.Helper()

	return commandContext{
		cacheRoot: t.TempDir(),
		env: append(
			os.Environ(),
			"GOCACHE="+t.TempDir(),
			"GOWORK=off",
		),
	}
}

func writeTestLinter(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	repoRoot := repoRoot(t)
	lintModuleDir := filepath.Join(repoRoot, "lint")
	goMod := fmt.Sprintf(`module example.com/byob-linter-test

go 1.26

require github.com/delinoio/byob/lint v0.0.0

replace github.com/delinoio/byob/lint => %s
`, filepath.ToSlash(lintModuleDir))
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	mainPath := filepath.Join(dir, "main.go")
	mainGo := `package main

import (
	"fmt"
	"os"
	"strings"

	byoblint "github.com/delinoio/byob/lint"
)

var _ = byoblint.Rule{
	Name: "example",
	Run: func(ctx byoblint.RuleContext, options any) byoblint.RuleListeners {
		return byoblint.RuleListeners{}
	},
}

func main() {
	fmt.Printf("args=%s\n", strings.Join(os.Args[1:], ","))
	if os.Getenv("BYOB_TEST_EXIT_7") == "1" {
		os.Exit(7)
	}
}
`
	if err := os.WriteFile(mainPath, []byte(mainGo), 0644); err != nil {
		t.Fatal(err)
	}
	return mainPath
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../../../.."))
}

func nonHostTarget() string {
	if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		return "darwin/arm64"
	}
	return "linux/amd64"
}
