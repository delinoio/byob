package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	byoblint "github.com/delinoio/byob/lint"
	bridge "github.com/microsoft/typescript-go/byobbridge"
)

type lintCommandIdentifier string

const (
	lintCommandBuild lintCommandIdentifier = "build"
	lintCommandRun   lintCommandIdentifier = "run"
)

type goTarget struct {
	GOOS   string
	GOARCH string
}

func hostTarget() goTarget {
	return goTarget{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH}
}

func parseGoTarget(raw string) (goTarget, error) {
	if raw == "" {
		return hostTarget(), nil
	}

	goos, goarch, ok := strings.Cut(raw, "/")
	if !ok || goos == "" || goarch == "" || strings.Contains(goarch, "/") {
		return goTarget{}, fmt.Errorf("target must use goos/goarch format")
	}
	return goTarget{GOOS: goos, GOARCH: goarch}, nil
}

func (target goTarget) String() string {
	return target.GOOS + "/" + target.GOARCH
}

func (target goTarget) isHost() bool {
	return target == hostTarget()
}

type lintBuildOptions struct {
	mainPath string
	outDir   string
	target   goTarget
	force    bool
}

type lintRunOptions struct {
	mainPath   string
	target     goTarget
	linterArgs []string
}

type lintMainPackage struct {
	mainPath   string
	packageDir string
	importPath string
	module     *goListModule
}

type lintBuildResult struct {
	binaryPath          string
	cacheKey            string
	mainPath            string
	packageDir          string
	packageImportPath   string
	target              goTarget
	built               bool
	bridgeVersion       string
	rslintCompatVersion string
}

type lintArtifactManifest struct {
	Tool                string `json:"tool"`
	Target              string `json:"target"`
	CacheKey            string `json:"cacheKey"`
	Binary              string `json:"binary"`
	SourceMain          string `json:"sourceMain"`
	Package             string `json:"package"`
	BYOBVersion         string `json:"byobVersion"`
	BridgeVersion       string `json:"bridgeVersion"`
	RslintCompatVersion string `json:"rslintCompatVersion"`
}

func executeLint(args []string, ctx commandContext) int {
	if len(args) == 0 {
		printLintUsage(ctx.stderr)
		return 2
	}

	switch lintCommandIdentifier(args[0]) {
	case "-h", "--help", "help":
		printLintUsage(ctx.stdout)
		return 0
	case lintCommandBuild:
		return executeLintBuild(args[1:], ctx)
	case lintCommandRun:
		return executeLintRun(args[1:], ctx)
	default:
		_, _ = fmt.Fprintf(ctx.stderr, "unknown lint command %q\n", args[0])
		printLintUsage(ctx.stderr)
		return 2
	}
}

func executeLintBuild(args []string, ctx commandContext) int {
	opts, ok := parseLintBuildOptions(args, ctx.stderr)
	if !ok {
		return 2
	}

	result, err := ensureLintBinary(ctx, opts)
	if err != nil {
		_, _ = fmt.Fprintf(ctx.stderr, "lint build failed: %v\n", err)
		return 1
	}

	if opts.outDir != "" {
		if err := exportLintArtifact(result, opts.outDir); err != nil {
			_, _ = fmt.Fprintf(ctx.stderr, "lint artifact export failed: %v\n", err)
			return 1
		}
		_, _ = fmt.Fprintf(ctx.stderr, "exported lint artifact: %s\n", opts.outDir)
	}

	return 0
}

func executeLintRun(args []string, ctx commandContext) int {
	opts, ok := parseLintRunOptions(args, ctx.stderr)
	if !ok {
		return 2
	}
	if !opts.target.isHost() {
		_, _ = fmt.Fprintf(ctx.stderr, "lint run target %s cannot execute on host %s\n", opts.target, hostTarget())
		return 2
	}

	result, err := ensureLintBinary(ctx, lintBuildOptions{
		mainPath: opts.mainPath,
		target:   opts.target,
	})
	if err != nil {
		_, _ = fmt.Fprintf(ctx.stderr, "lint run failed: %v\n", err)
		return 1
	}

	cmd := exec.Command(result.binaryPath, opts.linterArgs...)
	cmd.Stdin = ctx.stdin
	cmd.Stdout = ctx.stdout
	cmd.Stderr = ctx.stderr
	cmd.Env = ctx.env
	err = cmd.Run()
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	_, _ = fmt.Fprintf(ctx.stderr, "failed to execute lint binary %s: %v\n", result.binaryPath, err)
	return 1
}

func parseLintBuildOptions(args []string, stderr io.Writer) (lintBuildOptions, bool) {
	var opts lintBuildOptions
	var targetRaw string
	fs := flag.NewFlagSet("byob lint build", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.mainPath, "main", "", "path to the linter main.go file")
	fs.StringVar(&opts.outDir, "out", "", "directory to export the built artifact")
	fs.StringVar(&targetRaw, "target", "", "target goos/goarch")
	fs.BoolVar(&opts.force, "force", false, "rebuild even when a cached binary exists")
	fs.Usage = func() {
		_, _ = fmt.Fprintln(stderr, "usage: byob lint build --main <path> [--out <dir>] [--target <goos/goarch>] [--force]")
	}

	if err := fs.Parse(args); err != nil {
		return opts, false
	}
	if len(fs.Args()) != 0 {
		_, _ = fmt.Fprintln(stderr, "lint build accepts no positional arguments")
		fs.Usage()
		return opts, false
	}
	if opts.mainPath == "" {
		_, _ = fmt.Fprintln(stderr, "lint build requires --main")
		fs.Usage()
		return opts, false
	}

	target, err := parseGoTarget(targetRaw)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "invalid --target: %v\n", err)
		return opts, false
	}
	opts.target = target
	return opts, true
}

func parseLintRunOptions(args []string, stderr io.Writer) (lintRunOptions, bool) {
	var opts lintRunOptions
	var targetRaw string
	fs := flag.NewFlagSet("byob lint run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.mainPath, "main", "", "path to the linter main.go file")
	fs.StringVar(&targetRaw, "target", "", "target goos/goarch")
	fs.Usage = func() {
		_, _ = fmt.Fprintln(stderr, "usage: byob lint run --main <path> [--target <goos/goarch>] [--] <args...>")
	}

	if err := fs.Parse(args); err != nil {
		return opts, false
	}
	if opts.mainPath == "" {
		_, _ = fmt.Fprintln(stderr, "lint run requires --main")
		fs.Usage()
		return opts, false
	}

	target, err := parseGoTarget(targetRaw)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "invalid --target: %v\n", err)
		return opts, false
	}
	opts.target = target
	opts.linterArgs = fs.Args()
	return opts, true
}

func ensureLintBinary(ctx commandContext, opts lintBuildOptions) (lintBuildResult, error) {
	mainPkg, err := resolveLintMain(ctx, opts.mainPath, opts.target)
	if err != nil {
		return lintBuildResult{}, err
	}

	cacheKey, err := computeLintCacheKey(ctx, mainPkg, opts.target)
	if err != nil {
		return lintBuildResult{}, err
	}

	cacheRoot, err := lintCacheRoot(ctx)
	if err != nil {
		return lintBuildResult{}, err
	}

	cacheDir := filepath.Join(cacheRoot, cacheKey)
	binaryPath := filepath.Join(cacheDir, lintBinaryName(opts.target))
	result := lintBuildResult{
		binaryPath:          binaryPath,
		cacheKey:            cacheKey,
		mainPath:            mainPkg.mainPath,
		packageDir:          mainPkg.packageDir,
		packageImportPath:   mainPkg.importPath,
		target:              opts.target,
		bridgeVersion:       bridge.LinkedVersion(),
		rslintCompatVersion: byoblint.RslintCompatVersion,
	}

	if !opts.force && fileExists(binaryPath) {
		_, _ = fmt.Fprintf(ctx.stderr, "using cached lint binary: %s\n", binaryPath)
		return result, nil
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return result, err
	}

	tmp, err := os.CreateTemp(cacheDir, ".byob-lint-build-*")
	if err != nil {
		return result, err
	}
	tmpPath := tmp.Name()
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return result, err
	}
	_ = os.Remove(tmpPath)

	buildCmd := exec.Command("go", "build", "-trimpath", "-o", tmpPath, ".")
	buildCmd.Dir = mainPkg.packageDir
	buildCmd.Env = withEnv(ctx.env, "GOOS", opts.target.GOOS, "GOARCH", opts.target.GOARCH)
	buildCmd.Stdout = ctx.stderr
	buildCmd.Stderr = ctx.stderr
	if err := buildCmd.Run(); err != nil {
		_ = os.Remove(tmpPath)
		return result, err
	}

	if err := os.Chmod(tmpPath, 0755); err != nil {
		_ = os.Remove(tmpPath)
		return result, err
	}
	if fileExists(binaryPath) {
		_ = os.Remove(binaryPath)
	}
	if err := os.Rename(tmpPath, binaryPath); err != nil {
		_ = os.Remove(tmpPath)
		return result, err
	}

	result.built = true
	_, _ = fmt.Fprintf(ctx.stderr, "built lint binary: %s\n", binaryPath)
	return result, nil
}

func resolveLintMain(ctx commandContext, mainPath string, target goTarget) (lintMainPackage, error) {
	absMain, err := filepath.Abs(mainPath)
	if err != nil {
		return lintMainPackage{}, err
	}
	absMain = filepath.Clean(absMain)

	info, err := os.Stat(absMain)
	if err != nil {
		return lintMainPackage{}, err
	}
	if info.IsDir() {
		return lintMainPackage{}, fmt.Errorf("--main must point to a Go file, got directory %s", absMain)
	}
	if filepath.Ext(absMain) != ".go" {
		return lintMainPackage{}, fmt.Errorf("--main must point to a .go file")
	}

	pkg, err := goListPackage(ctx, filepath.Dir(absMain), target)
	if err != nil {
		return lintMainPackage{}, err
	}
	if pkg.Name != "main" {
		return lintMainPackage{}, fmt.Errorf("linter package must be package main, got package %s", pkg.Name)
	}

	return lintMainPackage{
		mainPath:   absMain,
		packageDir: pkg.Dir,
		importPath: pkg.ImportPath,
		module:     pkg.Module,
	}, nil
}

func computeLintCacheKey(ctx commandContext, mainPkg lintMainPackage, target goTarget) (string, error) {
	envInfo, err := goEnv(ctx, mainPkg.packageDir, target)
	if err != nil {
		return "", err
	}
	deps, err := goListDeps(ctx, mainPkg.packageDir, target)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hashString(hash, "byob-lint-cache-v1")
	hashString(hash, target.String())
	hashString(hash, byobVersion)
	hashString(hash, bridge.LinkedVersion())
	hashString(hash, byoblint.RuntimeVersion)
	hashString(hash, byoblint.RslintCompatVersion)
	hashString(hash, mainPkg.mainPath)
	hashString(hash, mainPkg.packageDir)
	hashString(hash, mainPkg.importPath)
	if mainPkg.module != nil {
		hashString(hash, mainPkg.module.Path)
		hashString(hash, mainPkg.module.Version)
		hashString(hash, mainPkg.module.Dir)
	}

	for _, path := range dependencyInputFiles(envInfo) {
		hashFile(hash, path)
	}

	sort.Slice(deps, func(i int, j int) bool {
		left := deps[i].ImportPath + "\x00" + deps[i].Dir
		right := deps[j].ImportPath + "\x00" + deps[j].Dir
		return left < right
	})
	for _, dep := range deps {
		if dep.Standard || dep.Dir == "" || isExternalGoPackageDir(dep.Dir, envInfo) {
			continue
		}
		hashString(hash, dep.ImportPath)
		hashString(hash, dep.Dir)
		for _, path := range dep.sourceFilePaths() {
			hashFile(hash, filepath.Join(dep.Dir, path))
		}
	}

	return hex.EncodeToString(hash.Sum(nil))[:32], nil
}

func lintCacheRoot(ctx commandContext) (string, error) {
	if ctx.cacheRoot != "" {
		return filepath.Join(ctx.cacheRoot, "byob", string(byoblint.ToolKindLint)), nil
	}
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheRoot, "byob", string(byoblint.ToolKindLint)), nil
}

func exportLintArtifact(result lintBuildResult, outDir string) error {
	absOut, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(absOut, 0755); err != nil {
		return err
	}

	binaryName := filepath.Base(result.binaryPath)
	exportedBinary := filepath.Join(absOut, binaryName)
	if err := copyFile(result.binaryPath, exportedBinary, 0755); err != nil {
		return err
	}

	manifest := lintArtifactManifest{
		Tool:                string(byoblint.ToolKindLint),
		Target:              result.target.String(),
		CacheKey:            result.cacheKey,
		Binary:              binaryName,
		SourceMain:          result.mainPath,
		Package:             result.packageImportPath,
		BYOBVersion:         byobVersion,
		BridgeVersion:       result.bridgeVersion,
		RslintCompatVersion: result.rslintCompatVersion,
	}
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	manifestBytes = append(manifestBytes, '\n')
	return os.WriteFile(filepath.Join(absOut, "byob-lint-artifact.json"), manifestBytes, 0644)
}

func printLintUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "usage: byob lint <command>")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "commands:")
	_, _ = fmt.Fprintln(w, "  build  Build and cache a user linter")
	_, _ = fmt.Fprintln(w, "  run    Build or reuse a cached linter, then execute it")
}

func lintBinaryName(target goTarget) string {
	name := "byob-lint-" + target.GOOS + "-" + target.GOARCH
	if target.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyFile(src string, dst string, mode os.FileMode) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func withEnv(env []string, pairs ...string) []string {
	next := append([]string{}, env...)
	for i := 0; i+1 < len(pairs); i += 2 {
		key := pairs[i]
		value := pairs[i+1]
		prefix := key + "="
		replaced := false
		for idx, existing := range next {
			if strings.HasPrefix(existing, prefix) {
				next[idx] = prefix + value
				replaced = true
				break
			}
		}
		if !replaced {
			next = append(next, prefix+value)
		}
	}
	return next
}
