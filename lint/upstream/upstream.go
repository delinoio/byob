package upstream

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	rslintbridge "github.com/web-infra-dev/rslint/byobbridge"
)

type repeatedFlag []string

func (flagValue *repeatedFlag) String() string {
	return strings.Join(*flagValue, ", ")
}

func (flagValue *repeatedFlag) Set(value string) error {
	*flagValue = append(*flagValue, value)
	return nil
}

func RunAllRulesCLI(args []string) int {
	return RunAllRulesCLIWithIO(args, os.Stdin, os.Stdout, os.Stderr)
}

func RunAllRulesCLIWithIO(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	var (
		allRules    bool
		configPath  string
		formatRaw   string
		fix         bool
		typeCheck   bool
		quiet       bool
		maxWarnings int
		ruleFlags   repeatedFlag
	)

	fs := flag.NewFlagSet("byob upstream lint", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.BoolVar(&allRules, "all-rules", false, "enable every pinned upstream rslint rule")
	fs.StringVar(&configPath, "config", "", "path to an rslint JSON/JSONC config")
	fs.StringVar(&formatRaw, "format", string(rslintbridge.OutputFormatDefault), "output format: default | jsonline | github")
	fs.BoolVar(&fix, "fix", false, "apply safe autofixes")
	fs.BoolVar(&typeCheck, "type-check", false, "enable TypeScript type-aware rules and semantic diagnostics")
	fs.BoolVar(&quiet, "quiet", false, "report errors only")
	fs.IntVar(&maxWarnings, "max-warnings", -1, "number of warnings that trigger a nonzero exit code")
	fs.Var(&ruleFlags, "rule", "rule override, e.g. 'no-console: off' (repeatable)")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: byob-upstream-lint [--all-rules] [--type-check] [--format default|jsonline|github] [--] <files-or-dirs...>")
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	_ = allRules

	format := rslintbridge.OutputFormat(formatRaw)
	switch format {
	case rslintbridge.OutputFormatDefault, rslintbridge.OutputFormatJSONLine, rslintbridge.OutputFormatGitHub:
	default:
		fmt.Fprintf(stderr, "invalid --format %q\n", formatRaw)
		return 2
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "failed to resolve current directory: %v\n", err)
		return 1
	}

	return rslintbridge.RunAllRules(rslintbridge.RunOptions{
		CWD:         cwd,
		Files:       fs.Args(),
		ConfigPath:  configPath,
		Format:      format,
		Fix:         fix,
		TypeCheck:   typeCheck,
		Quiet:       quiet,
		MaxWarnings: maxWarnings,
		RuleFlags:   []string(ruleFlags),
		Stdin:       stdin,
		Stdout:      stdout,
		Stderr:      stderr,
	})
}
