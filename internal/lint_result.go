package internal

import "strings"

type LintResult struct {
	Valid      bool                   `json:"valid"`
	Operations []*LintOperationResult `json:"operations"`
}

type LintOperationResult struct {
	Valid       bool   `json:"valid"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	StartLine   int32  `json:"startLine"`
	EndLine     int32  `json:"endLine"`
	Description string `json:"description"`
	HowToFix    string `json:"howToFix"`
}

type LintOperationResults []*LintOperationResult

func (r LintOperationResults) Len() int           { return len(r) }
func (r LintOperationResults) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r LintOperationResults) Less(i, j int) bool { return strings.Compare(r[i].Path, r[j].Path) < 0 }
