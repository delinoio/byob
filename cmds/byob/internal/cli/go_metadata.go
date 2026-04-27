package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type goListModule struct {
	Path    string
	Version string
	Dir     string
}

type goListPackageInfo struct {
	ImportPath   string
	Name         string
	Dir          string
	Standard     bool
	Module       *goListModule
	GoFiles      []string
	CgoFiles     []string
	CFiles       []string
	CXXFiles     []string
	MFiles       []string
	HFiles       []string
	FFiles       []string
	SFiles       []string
	SwigFiles    []string
	SwigCXXFiles []string
	SysoFiles    []string
	EmbedFiles   []string
}

func (pkg goListPackageInfo) sourceFilePaths() []string {
	var paths []string
	paths = append(paths, pkg.GoFiles...)
	paths = append(paths, pkg.CgoFiles...)
	paths = append(paths, pkg.CFiles...)
	paths = append(paths, pkg.CXXFiles...)
	paths = append(paths, pkg.MFiles...)
	paths = append(paths, pkg.HFiles...)
	paths = append(paths, pkg.FFiles...)
	paths = append(paths, pkg.SFiles...)
	paths = append(paths, pkg.SwigFiles...)
	paths = append(paths, pkg.SwigCXXFiles...)
	paths = append(paths, pkg.SysoFiles...)
	paths = append(paths, pkg.EmbedFiles...)
	sort.Strings(paths)
	return compactStrings(paths)
}

type goEnvInfo struct {
	GOMOD      string
	GOWORK     string
	GOMODCACHE string
	GOROOT     string
}

func goListPackage(ctx commandContext, dir string, target goTarget) (goListPackageInfo, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("go", "list", "-json", ".")
	cmd.Dir = dir
	cmd.Env = withEnv(ctx.env, "GOOS", target.GOOS, "GOARCH", target.GOARCH)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return goListPackageInfo{}, fmt.Errorf("go list failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	var pkg goListPackageInfo
	if err := json.Unmarshal(stdout.Bytes(), &pkg); err != nil {
		return goListPackageInfo{}, err
	}
	return pkg, nil
}

func goListDeps(ctx commandContext, dir string, target goTarget) ([]goListPackageInfo, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("go", "list", "-deps", "-json", ".")
	cmd.Dir = dir
	cmd.Env = withEnv(ctx.env, "GOOS", target.GOOS, "GOARCH", target.GOARCH)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go list deps failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	decoder := json.NewDecoder(&stdout)
	var pkgs []goListPackageInfo
	for {
		var pkg goListPackageInfo
		if err := decoder.Decode(&pkg); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

func goEnv(ctx commandContext, dir string, target goTarget) (goEnvInfo, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("go", "env", "-json", "GOMOD", "GOWORK", "GOMODCACHE", "GOROOT")
	cmd.Dir = dir
	cmd.Env = withEnv(ctx.env, "GOOS", target.GOOS, "GOARCH", target.GOARCH)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return goEnvInfo{}, fmt.Errorf("go env failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	var envInfo goEnvInfo
	if err := json.Unmarshal(stdout.Bytes(), &envInfo); err != nil {
		return goEnvInfo{}, err
	}
	return envInfo, nil
}

func dependencyInputFiles(envInfo goEnvInfo) []string {
	var paths []string
	if isUsableGoEnvPath(envInfo.GOMOD) {
		paths = append(paths, envInfo.GOMOD, filepath.Join(filepath.Dir(envInfo.GOMOD), "go.sum"))
	}
	if isUsableGoEnvPath(envInfo.GOWORK) {
		paths = append(paths, envInfo.GOWORK, envInfo.GOWORK+".sum")
	}
	sort.Strings(paths)
	return compactStrings(paths)
}

func isUsableGoEnvPath(path string) bool {
	return path != "" && path != "off" && path != os.DevNull
}

func isExternalGoPackageDir(dir string, envInfo goEnvInfo) bool {
	cleanDir := filepath.Clean(dir)
	for _, root := range []string{envInfo.GOMODCACHE, envInfo.GOROOT} {
		if root == "" {
			continue
		}
		cleanRoot := filepath.Clean(root)
		if cleanDir == cleanRoot || strings.HasPrefix(cleanDir, cleanRoot+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func hashString(h hash.Hash, value string) {
	_, _ = h.Write([]byte(value))
	_, _ = h.Write([]byte{0})
}

func hashFile(h hash.Hash, path string) {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		hashString(h, "missing:"+path)
		return
	}
	hashString(h, "file:"+path)
	contents, err := os.ReadFile(path)
	if err != nil {
		hashString(h, "unreadable")
		return
	}
	_, _ = h.Write(contents)
	_, _ = h.Write([]byte{0})
}

func compactStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	compact := values[:0]
	var previous string
	for index, value := range values {
		if index == 0 || value != previous {
			compact = append(compact, value)
			previous = value
		}
	}
	return compact
}
