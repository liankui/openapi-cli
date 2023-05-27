package pkg

import "strings"

type LintOperationResults []*LintOperationResult

func (r LintOperationResults) Len() int {
	return len(r)
}

func (r LintOperationResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r LintOperationResults) Less(i, j int) bool {
	return strings.Compare(r[i].Path, r[j].Path) < 0
}
