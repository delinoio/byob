package byobbridge

import (
	"slices"
	"testing"
)

func TestInfoReportsPinnedUpstream(t *testing.T) {
	info := Info()
	if info.Module != "github.com/web-infra-dev/rslint" {
		t.Fatalf("unexpected module: %s", info.Module)
	}
	if info.Commit != "b986707ef58537329229645b77d35d93e3234ca6" {
		t.Fatalf("unexpected commit: %s", info.Commit)
	}
	if info.RuleCount != len(AllRuleNames()) {
		t.Fatalf("unexpected rule count: %d", info.RuleCount)
	}
}

func TestAllRuleNamesIncludesRepresentativeUpstreamRules(t *testing.T) {
	names := AllRuleNames()
	if len(names) != 282 {
		t.Fatalf("unexpected upstream rule count: %d", len(names))
	}

	for _, name := range []string{
		"no-console",
		"@typescript-eslint/no-unused-vars",
		"import/no-duplicates",
		"jest/no-focused-tests",
		"promise/param-names",
		"react/jsx-key",
		"react-hooks/rules-of-hooks",
		"unicorn/filename-case",
	} {
		if !slices.Contains(names, name) {
			t.Fatalf("expected rule %q in upstream registry", name)
		}
	}
}
