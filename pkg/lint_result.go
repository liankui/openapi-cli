package pkg

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
