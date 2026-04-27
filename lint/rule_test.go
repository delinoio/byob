package lint

import "testing"

type fakeNode struct {
	kind      NodeKind
	textRange TextRange
}

func (node fakeNode) Kind() NodeKind {
	return node.kind
}

func (node fakeNode) Range() TextRange {
	return node.textRange
}

func TestCreateRulePrefixesTypeScriptESLintName(t *testing.T) {
	rule := CreateRule(Rule{Name: "no-unused-vars", RequiresTypeInfo: true})
	if rule.Name != "@typescript-eslint/no-unused-vars" {
		t.Fatalf("unexpected rule name: %s", rule.Name)
	}
	if !rule.RequiresTypeInfo {
		t.Fatal("expected RequiresTypeInfo to be preserved")
	}
}

func TestSeverityParsing(t *testing.T) {
	if ParseSeverity("warning") != SeverityWarning {
		t.Fatal("expected warning to parse as SeverityWarning")
	}
	if ParseSeverity("off").Int() != 0 {
		t.Fatal("expected off severity int to be 0")
	}
	if ParseSeverity("unknown") != SeverityError {
		t.Fatal("expected unknown severity to default to SeverityError")
	}
}

func TestRuleFixHelpers(t *testing.T) {
	node := fakeNode{kind: 1, textRange: NewTextRange(4, 9)}

	before := RuleFixInsertBefore(node, "x")
	if before.Range.Pos() != 4 || before.Range.End() != 4 {
		t.Fatalf("unexpected insert-before range: %#v", before.Range)
	}

	after := RuleFixInsertAfter(node, "x")
	if after.Range.Pos() != 9 || after.Range.End() != 9 {
		t.Fatalf("unexpected insert-after range: %#v", after.Range)
	}

	remove := RuleFixRemove(node)
	if remove.Text != "" || remove.Range.Pos() != 4 || remove.Range.End() != 9 {
		t.Fatalf("unexpected remove fix: %#v", remove)
	}
}

func TestRuleDiagnosticFixes(t *testing.T) {
	if fixes := (RuleDiagnostic{}).Fixes(); len(fixes) != 0 {
		t.Fatalf("expected empty fixes for nil FixesPtr, got=%d", len(fixes))
	}

	expected := []RuleFix{{Text: "x", Range: NewTextRange(1, 2)}}
	diagnostic := RuleDiagnostic{FixesPtr: &expected}
	if fixes := diagnostic.Fixes(); len(fixes) != 1 || fixes[0].Text != "x" {
		t.Fatalf("unexpected fixes: %#v", fixes)
	}
}

func TestRuleShapeCompiles(t *testing.T) {
	rule := Rule{
		Name: "example",
		Run: func(ctx RuleContext, options any) RuleListeners {
			return RuleListeners{
				1: func(node Node) {
					ctx.ReportNode(node, RuleMessage{Id: "example", Description: "example diagnostic"})
				},
			}
		},
	}
	if rule.Run == nil {
		t.Fatal("expected rule run function")
	}
}
