package pkg

import (
	"context"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/daveshanley/vacuum/rulesets"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	LintRules = &rulesets.RuleSet{
		Rules: map[string]*model.Rule{
			"operation-parameters": rulesets.GetOperationParametersRule(),
		},
	}
)

type Lint struct{}

func NewLint() *Lint {
	return &Lint{}
}

func (l *Lint) Action(c *cli.Context) error {
	slog.Info("api lint", "file", c.Args().First())

	spec, err := os.ReadFile(c.Args().First())
	if err != nil {
		slog.Error("failed to read file", "error", err, "path", os.Args[1:])
		return nil
	}

	lintResult, err := OpenapiLint(c.Context, spec)
	if err != nil {
		slog.Warn("[Action] openapi lint error", "error", err)
		return nil
	}

	for _, o := range lintResult.Operations {
		if !o.Valid {
			slog.Info("violation", "result", o)
		}
	}

	slog.Info("api lint finished", "file", c.Args().First())

	return nil
}

func OpenapiLint(ctx context.Context, spec []byte) (*LintResult, error) {
	result := motor.ApplyRulesToRuleSet(&motor.RuleSetExecution{
		RuleSet:         LintRules,
		Spec:            spec,
		CustomFunctions: map[string]model.RuleFunction{},
	})

	if result.Index == nil {
		return &LintResult{Valid: false}, nil
	}

	if len(result.Errors) > 0 {
		return &LintResult{Valid: false}, errors.Errorf("apply rule get errors, errors: %v", result.Errors)
	}

	operations := result.Index.GetAllPaths()

	lintResult := &LintResult{
		Operations: make([]*LintOperationResult, 0, len(operations)),
		Valid:      true,
	}

	for path, operation := range operations {
		for method := range operation {
			operationResult := &LintOperationResult{
				Path:   path,
				Method: method,
				Valid:  true,
			}

			operationPath := strings.Join([]string{"$.paths", path, method, "parameters"}, ".")
			for _, _result := range result.Results {
				if _result.Path != operationPath {
					continue
				}

				operationResult.Valid = false
				operationResult.Description = _result.Rule.Description
				operationResult.HowToFix = _result.Rule.HowToFix
				if _result.StartNode != nil {
					operationResult.StartLine = int32(_result.StartNode.Line)
				}
				if _result.EndNode != nil {
					operationResult.EndLine = int32(_result.EndNode.Line)
				}

				lintResult.Valid = false
				break
			}

			lintResult.Operations = append(lintResult.Operations, operationResult)
		}
	}

	sort.Sort(LintOperationResults(lintResult.Operations))

	return lintResult, nil
}
