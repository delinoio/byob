package byobbridge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/microsoft/typescript-go/shim/compiler"
	"github.com/microsoft/typescript-go/shim/core"
	"github.com/microsoft/typescript-go/shim/scanner"
	"github.com/microsoft/typescript-go/shim/tspath"
	"github.com/microsoft/typescript-go/shim/vfs"
	"github.com/microsoft/typescript-go/shim/vfs/osvfs"
	"github.com/web-infra-dev/rslint/internal/config"
	"github.com/web-infra-dev/rslint/internal/linter"
	"github.com/web-infra-dev/rslint/internal/rule"
	"github.com/web-infra-dev/rslint/internal/utils"
)

const (
	upstreamModule = "github.com/web-infra-dev/rslint"
	upstreamCommit = "b986707ef58537329229645b77d35d93e3234ca6"
)

type OutputFormat string

const (
	OutputFormatDefault  OutputFormat = "default"
	OutputFormatJSONLine OutputFormat = "jsonline"
	OutputFormatGitHub   OutputFormat = "github"
)

type LinkInfo struct {
	Module    string
	Commit    string
	RuleCount int
}

type RunOptions struct {
	CWD         string
	Files       []string
	ConfigPath  string
	Format      OutputFormat
	Fix         bool
	TypeCheck   bool
	Quiet       bool
	MaxWarnings int
	RuleFlags   []string
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
}

type diagnosticJSON struct {
	RuleName string        `json:"ruleName"`
	Message  string        `json:"message"`
	FilePath string        `json:"filePath"`
	Range    diagnosticLoc `json:"range"`
	Severity string        `json:"severity"`
}

type diagnosticLoc struct {
	Start position `json:"start"`
	End   position `json:"end"`
}

type position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func Info() LinkInfo {
	return LinkInfo{
		Module:    upstreamModule,
		Commit:    upstreamCommit,
		RuleCount: len(AllRuleNames()),
	}
}

func UpstreamCommit() string {
	return upstreamCommit
}

func AllRuleNames() []string {
	config.RegisterAllRules()
	rules := config.GlobalRuleRegistry.GetAllRules()
	names := make([]string, 0, len(rules))
	for name := range rules {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func RunAllRules(options RunOptions) int {
	runner := allRulesRunner{options: normalizeRunOptions(options)}
	return runner.run()
}

type allRulesRunner struct {
	options            RunOptions
	fsys               vfs.FS
	cwd                string
	comparePathOptions tspath.ComparePathsOptions
}

func normalizeRunOptions(options RunOptions) RunOptions {
	if options.Format == "" {
		options.Format = OutputFormatDefault
	}
	if options.Stdin == nil {
		options.Stdin = os.Stdin
	}
	if options.Stdout == nil {
		options.Stdout = os.Stdout
	}
	if options.Stderr == nil {
		options.Stderr = os.Stderr
	}
	return options
}

func (runner *allRulesRunner) run() int {
	config.RegisterAllRules()

	cwd, err := runner.resolveCWD()
	if err != nil {
		fmt.Fprintf(runner.options.Stderr, "error: %v\n", err)
		return 1
	}
	runner.cwd = cwd
	runner.fsys = osvfs.FS()
	runner.comparePathOptions = tspath.ComparePathsOptions{
		CurrentDirectory:          cwd,
		UseCaseSensitiveFileNames: runner.fsys.UseCaseSensitiveFileNames(),
	}

	rslintConfig, configDir, err := runner.loadConfig()
	if err != nil {
		fmt.Fprintf(runner.options.Stderr, "error: %v\n", err)
		return 1
	}

	allowFiles, allowDirs, err := runner.resolveInputs()
	if err != nil {
		fmt.Fprintf(runner.options.Stderr, "error: %v\n", err)
		return 1
	}

	collect := func() ([]rule.RuleDiagnostic, *linter.LintResult, error) {
		programs, typeInfoFiles, err := runner.createPrograms(rslintConfig, configDir, allowFiles, allowDirs)
		if err != nil {
			return nil, nil, err
		}
		if len(programs) == 0 {
			return nil, &linter.LintResult{ExecutedRules: map[string]struct{}{}}, nil
		}

		var diagnostics []rule.RuleDiagnostic
		var diagnosticsMu sync.Mutex
		result, err := linter.RunLinter(
			programs,
			true,
			allowFiles,
			allowDirs,
			utils.ExcludePaths,
			func(sourceFile *ast.SourceFile) []linter.ConfiguredRule {
				return config.GlobalRuleRegistry.GetActiveRulesForFile(rslintConfig, sourceFile.FileName(), configDir, false, typeInfoFiles)
			},
			runner.options.TypeCheck,
			func(diagnostic rule.RuleDiagnostic) {
				diagnosticsMu.Lock()
				diagnostics = append(diagnostics, diagnostic)
				diagnosticsMu.Unlock()
			},
			typeInfoFiles,
			runner.fileFilters(len(programs), rslintConfig, configDir),
		)
		return diagnostics, result, err
	}

	diagnostics, result, err := collect()
	if err != nil {
		fmt.Fprintf(runner.options.Stderr, "error: %v\n", err)
		return 1
	}

	fixedCount := 0
	if runner.options.Fix && len(diagnostics) > 0 {
		fixedCount = runner.applyFixes(diagnostics)
		if fixedCount > 0 {
			diagnostics, result, err = collect()
			if err != nil {
				fmt.Fprintf(runner.options.Stderr, "error: %v\n", err)
				return 1
			}
		}
	}

	return runner.finish(diagnostics, result, fixedCount)
}

func (runner *allRulesRunner) resolveCWD() (string, error) {
	if runner.options.CWD == "" {
		return os.Getwd()
	}
	return filepath.Abs(runner.options.CWD)
}

func (runner *allRulesRunner) loadConfig() (config.RslintConfig, string, error) {
	if runner.options.ConfigPath != "" {
		loader := config.NewConfigLoader(runner.fsys, runner.cwd)
		rslintConfig, configDir, err := loader.LoadRslintConfig(runner.options.ConfigPath)
		if err != nil {
			return nil, "", err
		}
		return runner.applyRuleFlags(rslintConfig), configDir, nil
	}

	rules := make(config.Rules)
	for _, name := range AllRuleNames() {
		rules[name] = "error"
	}

	entry := config.ConfigEntry{
		LanguageOptions: &config.LanguageOptions{},
		Plugins:         allPluginNames(),
		Rules:           rules,
	}
	if runner.options.TypeCheck {
		if tsconfigPath := runner.findTsConfig(); tsconfigPath != "" {
			entry.LanguageOptions.ParserOptions = &config.ParserOptions{
				Project: config.ProjectPaths{tsconfigPath},
			}
		}
	}

	return runner.applyRuleFlags(config.RslintConfig{entry}), runner.cwd, nil
}

func (runner *allRulesRunner) applyRuleFlags(rslintConfig config.RslintConfig) config.RslintConfig {
	if len(runner.options.RuleFlags) == 0 {
		return rslintConfig
	}
	entry, err := config.BuildCLIRuleEntry(runner.options.RuleFlags)
	if err != nil {
		fmt.Fprintf(runner.options.Stderr, "warning: ignoring invalid rule override: %v\n", err)
		return rslintConfig
	}
	if entry == nil {
		return rslintConfig
	}
	return append(rslintConfig, *entry)
}

func allPluginNames() []string {
	plugins := make([]string, 0, len(config.KnownPlugins))
	for _, plugin := range config.KnownPlugins {
		plugins = append(plugins, plugin.RulePrefix)
	}
	sort.Strings(plugins)
	return plugins
}

func (runner *allRulesRunner) findTsConfig() string {
	candidates := append([]string{}, runner.options.Files...)
	candidates = append(candidates, ".")
	for _, candidate := range candidates {
		abs, err := runner.absPath(candidate)
		if err != nil {
			continue
		}
		info, err := os.Stat(abs)
		if err == nil && !info.IsDir() {
			abs = filepath.Dir(abs)
		}
		for {
			tsconfig := filepath.Join(abs, "tsconfig.json")
			if fileExists(tsconfig) {
				return tspath.NormalizePath(tsconfig)
			}
			if filepath.Clean(abs) == filepath.Clean(runner.cwd) {
				break
			}
			parent := filepath.Dir(abs)
			if parent == abs {
				break
			}
			abs = parent
		}
	}
	return ""
}

func (runner *allRulesRunner) resolveInputs() ([]string, []string, error) {
	if len(runner.options.Files) == 0 {
		return nil, []string{tspath.NormalizePath(runner.cwd)}, nil
	}

	seenFiles := map[string]struct{}{}
	seenDirs := map[string]struct{}{}
	var files []string
	var dirs []string
	for _, raw := range runner.options.Files {
		abs, err := runner.absPath(raw)
		if err != nil {
			return nil, nil, err
		}
		info, err := os.Stat(abs)
		if err != nil {
			return nil, nil, err
		}
		normalized := tspath.NormalizePath(abs)
		if info.IsDir() {
			if _, ok := seenDirs[normalized]; !ok {
				seenDirs[normalized] = struct{}{}
				dirs = append(dirs, normalized)
			}
			continue
		}
		if _, ok := seenFiles[normalized]; !ok {
			seenFiles[normalized] = struct{}{}
			files = append(files, normalized)
		}
	}
	sort.Strings(files)
	sort.Strings(dirs)
	return files, dirs, nil
}

func (runner *allRulesRunner) absPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Abs(filepath.Join(runner.cwd, path))
}

func (runner *allRulesRunner) createPrograms(rslintConfig config.RslintConfig, configDir string, allowFiles []string, allowDirs []string) ([]*compiler.Program, map[string]struct{}, error) {
	host := utils.CreateCompilerHost(configDir, runner.fsys)
	var programs []*compiler.Program
	var typeInfoFiles map[string]struct{}

	tsConfigs, err := config.ResolveTsConfigPaths(rslintConfig, configDir, runner.fsys)
	if err != nil {
		return nil, nil, err
	}
	if len(tsConfigs) > 0 {
		for _, tsconfig := range tsConfigs {
			program, err := utils.CreateProgram(true, runner.fsys, configDir, tsconfig, host)
			if err != nil {
				return nil, nil, err
			}
			programs = append(programs, program)
		}
		typeInfoFiles = utils.CollectProgramFiles(programs, runner.fsys)
	}

	if len(programs) == 0 {
		rootFiles := runner.discoverRootFiles(configDir, allowFiles, allowDirs)
		if len(rootFiles) == 0 {
			return nil, typeInfoFiles, nil
		}
		program, err := utils.CreateProgramFromOptionsLenient(true, &core.CompilerOptions{
			Target:  core.ScriptTargetESNext,
			Module:  core.ModuleKindESNext,
			Jsx:     core.JsxEmitPreserve,
			AllowJs: core.TSTrue,
		}, rootFiles, host)
		if err != nil {
			return nil, nil, err
		}
		programs = append(programs, program)
		if runner.options.TypeCheck {
			typeInfoFiles = map[string]struct{}{}
		}
	}

	return programs, typeInfoFiles, nil
}

func (runner *allRulesRunner) discoverRootFiles(configDir string, allowFiles []string, allowDirs []string) []string {
	seen := make(map[string]struct{})
	var roots []string
	add := func(path string) {
		normalized := tspath.NormalizePath(path)
		if !isSupportedSourcePath(normalized) {
			return
		}
		if _, ok := seen[normalized]; ok {
			return
		}
		seen[normalized] = struct{}{}
		roots = append(roots, normalized)
	}

	for _, file := range allowFiles {
		add(file)
	}

	sourceExts := []string{".ts", ".tsx", ".js", ".jsx", ".mts", ".mjs", ".cts", ".cjs"}
	dirs := allowDirs
	if len(dirs) == 0 && len(allowFiles) == 0 {
		dirs = []string{configDir}
	}
	for _, dir := range dirs {
		for _, file := range vfs.ReadDirectory(runner.fsys, dir, configDir, sourceExts, utils.DefaultExcludeDirNames, []string{"**/*"}, nil) {
			add(file)
		}
	}
	sort.Strings(roots)
	return roots
}

func isSupportedSourcePath(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".ts", ".tsx", ".js", ".jsx", ".mts", ".mjs", ".cts", ".cjs":
		return true
	default:
		return false
	}
}

func (runner *allRulesRunner) fileFilters(count int, rslintConfig config.RslintConfig, configDir string) []func(string) bool {
	filters := make([]func(string) bool, count)
	for i := range filters {
		filters[i] = func(fileName string) bool {
			return !rslintConfig.IsFileIgnored(fileName, configDir)
		}
	}
	return filters
}

func (runner *allRulesRunner) applyFixes(diagnostics []rule.RuleDiagnostic) int {
	byFile := make(map[string][]rule.RuleDiagnostic)
	for _, diagnostic := range diagnostics {
		if len(diagnostic.Fixes()) == 0 || diagnostic.SourceFile == nil {
			continue
		}
		byFile[diagnostic.SourceFile.FileName()] = append(byFile[diagnostic.SourceFile.FileName()], diagnostic)
	}

	fixed := 0
	for fileName, fileDiagnostics := range byFile {
		contents, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(runner.options.Stderr, "warning: failed to read %s for fixes: %v\n", fileName, err)
			continue
		}
		fixedContent, unapplied, wasFixed := linter.ApplyRuleFixes(string(contents), fileDiagnostics)
		if !wasFixed {
			continue
		}
		if err := os.WriteFile(fileName, []byte(fixedContent), 0644); err != nil {
			fmt.Fprintf(runner.options.Stderr, "warning: failed to write %s for fixes: %v\n", fileName, err)
			continue
		}
		fixed += len(fileDiagnostics) - len(unapplied)
	}
	return fixed
}

func (runner *allRulesRunner) finish(diagnostics []rule.RuleDiagnostic, result *linter.LintResult, fixedCount int) int {
	sortDiagnostics(diagnostics)

	writer := bufio.NewWriter(runner.options.Stdout)
	defer writer.Flush()

	errorsCount := 0
	warningsCount := 0
	for _, diagnostic := range diagnostics {
		switch diagnostic.Severity {
		case rule.SeverityError:
			errorsCount++
		case rule.SeverityWarning:
			warningsCount++
		}
		if runner.options.Quiet && diagnostic.Severity != rule.SeverityError {
			continue
		}
		runner.printDiagnostic(writer, diagnostic)
	}

	if runner.options.Format == OutputFormatDefault {
		ruleCount := 0
		if result != nil {
			ruleCount = len(result.ExecutedRules)
		}
		fmt.Fprintf(writer, "\nFound %d error(s) and %d warning(s)", errorsCount, warningsCount)
		if fixedCount > 0 {
			fmt.Fprintf(writer, " (fixed %d issue(s))", fixedCount)
		}
		fmt.Fprintf(writer, " using %d rule(s)\n", ruleCount)
	}

	if errorsCount > 0 {
		return 1
	}
	if runner.options.MaxWarnings >= 0 && warningsCount > runner.options.MaxWarnings {
		return 1
	}
	return 0
}

func sortDiagnostics(diagnostics []rule.RuleDiagnostic) {
	sort.SliceStable(diagnostics, func(i int, j int) bool {
		left := diagnostics[i]
		right := diagnostics[j]
		leftFile := ""
		rightFile := ""
		if left.SourceFile != nil {
			leftFile = left.SourceFile.FileName()
		}
		if right.SourceFile != nil {
			rightFile = right.SourceFile.FileName()
		}
		if leftFile != rightFile {
			return leftFile < rightFile
		}
		if left.Range.Pos() != right.Range.Pos() {
			return left.Range.Pos() < right.Range.Pos()
		}
		return left.RuleName < right.RuleName
	})
}

func (runner *allRulesRunner) printDiagnostic(writer *bufio.Writer, diagnostic rule.RuleDiagnostic) {
	if diagnostic.SourceFile == nil {
		return
	}
	startLine, startColumn := scanner.GetECMALineAndUTF16CharacterOfPosition(diagnostic.SourceFile, diagnostic.Range.Pos())
	endLine, endColumn := scanner.GetECMALineAndUTF16CharacterOfPosition(diagnostic.SourceFile, diagnostic.Range.End())
	relativePath := tspath.ConvertToRelativePath(diagnostic.SourceFile.FileName(), runner.comparePathOptions)

	switch runner.options.Format {
	case OutputFormatJSONLine:
		payload := diagnosticJSON{
			RuleName: diagnostic.RuleName,
			Message:  diagnostic.Message.Description,
			FilePath: relativePath,
			Range: diagnosticLoc{
				Start: position{Line: startLine + 1, Column: int(startColumn) + 1},
				End:   position{Line: endLine + 1, Column: int(endColumn) + 1},
			},
			Severity: diagnostic.Severity.String(),
		}
		bytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Fprintf(writer, "{\"error\":%q}\n", err.Error())
			return
		}
		writer.Write(bytes)
		writer.WriteByte('\n')
	case OutputFormatGitHub:
		fmt.Fprintf(
			writer,
			"::%s file=%s,line=%d,endLine=%d,col=%d,endColumn=%d,title=%s::%s\n",
			githubSeverity(diagnostic.Severity),
			escapeGitHubProperty(relativePath),
			startLine+1,
			endLine+1,
			int(startColumn)+1,
			int(endColumn)+1,
			escapeGitHubProperty(diagnostic.RuleName),
			escapeGitHubData(diagnostic.Message.Description),
		)
	default:
		fmt.Fprintf(
			writer,
			"%s:%d:%d [%s] %s: %s\n",
			relativePath,
			startLine+1,
			int(startColumn)+1,
			diagnostic.Severity.String(),
			diagnostic.RuleName,
			diagnostic.Message.Description,
		)
	}
}

func githubSeverity(severity rule.DiagnosticSeverity) string {
	switch severity {
	case rule.SeverityError:
		return "error"
	case rule.SeverityWarning:
		return "warning"
	default:
		return "notice"
	}
}

func escapeGitHubData(value string) string {
	value = strings.ReplaceAll(value, "%", "%25")
	value = strings.ReplaceAll(value, "\r", "%0D")
	value = strings.ReplaceAll(value, "\n", "%0A")
	return value
}

func escapeGitHubProperty(value string) string {
	value = escapeGitHubData(value)
	value = strings.ReplaceAll(value, ":", "%3A")
	value = strings.ReplaceAll(value, ",", "%2C")
	return value
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
