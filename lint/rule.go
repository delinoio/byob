package lint

const (
	RuntimeVersion      = "0.1.0"
	RslintCompatVersion = "0.5.0"
)

type ToolKind string

const (
	ToolKindLint      ToolKind = "lint"
	ToolKindFormat    ToolKind = "fmt"
	ToolKindTransform ToolKind = "transform"
)

type DiagnosticSeverity int

const (
	SeverityError DiagnosticSeverity = iota
	SeverityWarning
	SeverityOff
)

func (severity DiagnosticSeverity) String() string {
	switch severity {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warn"
	case SeverityOff:
		return "off"
	default:
		return "error"
	}
}

func (severity DiagnosticSeverity) Int() int {
	switch severity {
	case SeverityError:
		return 1
	case SeverityWarning:
		return 2
	case SeverityOff:
		return 0
	default:
		return 0
	}
}

func ParseSeverity(level string) DiagnosticSeverity {
	switch level {
	case "error":
		return SeverityError
	case "warn", "warning":
		return SeverityWarning
	case "off":
		return SeverityOff
	default:
		return SeverityError
	}
}

type NodeKind int

const (
	listenerOnExitOffset                  NodeKind = 1000
	listenerOnAllowPatternOffset          NodeKind = 2000
	listenerOnAllowPatternOnExitOffset    NodeKind = 4000
	listenerOnNotAllowPatternOffset       NodeKind = 5000
	listenerOnNotAllowPatternOnExitOffset NodeKind = 6000
)

func ListenerOnExit(kind NodeKind) NodeKind {
	return kind + listenerOnExitOffset
}

func ListenerOnAllowPattern(kind NodeKind) NodeKind {
	return kind + listenerOnAllowPatternOffset
}

func ListenerOnNotAllowPattern(kind NodeKind) NodeKind {
	return kind + listenerOnAllowPatternOnExitOffset
}

type TextRange struct {
	PosValue int
	EndValue int
}

func NewTextRange(pos int, end int) TextRange {
	return TextRange{PosValue: pos, EndValue: end}
}

func (textRange TextRange) Pos() int {
	return textRange.PosValue
}

func (textRange TextRange) End() int {
	return textRange.EndValue
}

func (textRange TextRange) WithPos(pos int) TextRange {
	textRange.PosValue = pos
	return textRange
}

func (textRange TextRange) WithEnd(end int) TextRange {
	textRange.EndValue = end
	return textRange
}

type Node interface {
	Kind() NodeKind
	Range() TextRange
}

type SourceFile interface {
	FileName() string
	Text() string
}

type Program interface {
	SourceFiles() []SourceFile
}

type TypeChecker interface{}

type RuleListeners map[NodeKind]func(node Node)

type Rule struct {
	Name             string
	RequiresTypeInfo bool
	Run              func(ctx RuleContext, options any) RuleListeners
}

func CreateRule(rule Rule) Rule {
	return Rule{
		Name:             "@typescript-eslint/" + rule.Name,
		RequiresTypeInfo: rule.RequiresTypeInfo,
		Run:              rule.Run,
	}
}

func CreateTypeScriptRule(rule Rule) Rule {
	return CreateRule(rule)
}

type RuleMessage struct {
	Id          string
	Description string
}

type RuleFix struct {
	Text  string
	Range TextRange
}

func RuleFixInsertBefore(node Node, text string) RuleFix {
	nodeRange := node.Range()
	return RuleFix{Text: text, Range: nodeRange.WithEnd(nodeRange.Pos())}
}

func RuleFixInsertAfter(node Node, text string) RuleFix {
	nodeRange := node.Range()
	return RuleFix{Text: text, Range: nodeRange.WithPos(nodeRange.End())}
}

func RuleFixReplace(node Node, text string) RuleFix {
	return RuleFixReplaceRange(node.Range(), text)
}

func RuleFixReplaceRange(textRange TextRange, text string) RuleFix {
	return RuleFix{Text: text, Range: textRange}
}

func RuleFixRemove(node Node) RuleFix {
	return RuleFixReplace(node, "")
}

func RuleFixRemoveRange(textRange TextRange) RuleFix {
	return RuleFixReplaceRange(textRange, "")
}

type RuleSuggestion struct {
	Message  RuleMessage
	FixesArr []RuleFix
}

func (suggestion RuleSuggestion) Fixes() []RuleFix {
	return suggestion.FixesArr
}

type RuleDiagnostic struct {
	Range        TextRange
	RuleName     string
	Message      RuleMessage
	FixesPtr     *[]RuleFix
	Suggestions  *[]RuleSuggestion
	SourceFile   SourceFile
	Severity     DiagnosticSeverity
	PreFormatted bool
}

func (diagnostic RuleDiagnostic) Fixes() []RuleFix {
	if diagnostic.FixesPtr == nil {
		return []RuleFix{}
	}
	return *diagnostic.FixesPtr
}

type ReportRangeFunc func(textRange TextRange, message RuleMessage)
type ReportRangeWithFixesFunc func(textRange TextRange, message RuleMessage, fixes ...RuleFix)
type ReportRangeWithSuggestionsFunc func(textRange TextRange, message RuleMessage, suggestions ...RuleSuggestion)
type ReportNodeFunc func(node Node, message RuleMessage)
type ReportNodeWithFixesFunc func(node Node, message RuleMessage, fixes ...RuleFix)
type ReportNodeWithSuggestionsFunc func(node Node, message RuleMessage, suggestions ...RuleSuggestion)

type RuleContext struct {
	SourceFile                         SourceFile
	Settings                           map[string]any
	Program                            Program
	TypeChecker                        TypeChecker
	ReportRange                        ReportRangeFunc
	ReportRangeWithFixes               ReportRangeWithFixesFunc
	ReportRangeWithSuggestions         ReportRangeWithSuggestionsFunc
	ReportNode                         ReportNodeFunc
	ReportNodeWithFixes                ReportNodeWithFixesFunc
	ReportNodeWithSuggestions          ReportNodeWithSuggestionsFunc
	ReportNodeWithFixesAndSuggestions  func(node Node, message RuleMessage, fixes []RuleFix, suggestions []RuleSuggestion)
	ReportRangeWithFixesAndSuggestions func(textRange TextRange, message RuleMessage, fixes []RuleFix, suggestions []RuleSuggestion)
}

func ReportNodeWithFixesOrSuggestions(ctx RuleContext, node Node, fix bool, message RuleMessage, suggestionMessage RuleMessage, fixes ...RuleFix) {
	if fix {
		ctx.ReportNodeWithFixes(node, message, fixes...)
		return
	}
	ctx.ReportNodeWithSuggestions(node, message, RuleSuggestion{
		Message:  suggestionMessage,
		FixesArr: fixes,
	})
}
