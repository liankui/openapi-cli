package pkg

import (
	"context"
	"os"
	"sort"
	"strings"
	
	"github.com/chaos-io/chaos/core/logs"
	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/daveshanley/vacuum/rulesets"
	"github.com/urfave/cli/v2"
)

var (
	LintRules = &rulesets.RuleSet{
		Rules: map[string]*model.Rule{
			"operation-parameters": rulesets.GetOperationParametersRule(),
		},
	}
)

type LintConfig struct {
	InputFile string `json:"input_file"`
}

type Lint struct {
	config *LintConfig
}

func NewLint() *Lint {
	return &Lint{config: &LintConfig{}}
}

func (l *Lint) Action(c *cli.Context) error {
	spec, err := os.ReadFile(c.Args().First())
	if err != nil {
		logs.Errorw("failed to read file", "error", err, "path", os.Args[1:])
		return nil
	}
	
	lintResult := OpenapiLint(c.Context, spec)
	
	for _, o := range lintResult.Operations {
		if !o.Valid {
			logs.Infow("violation", "result", o)
		}
	}
	
	return nil
}

func OpenapiLint(ctx context.Context, spec []byte) *LintResult {
	result := motor.ApplyRulesToRuleSet(&motor.RuleSetExecution{
		RuleSet:         LintRules,
		Spec:            spec,
		CustomFunctions: map[string]model.RuleFunction{},
	})
	
	if result.Index == nil {
		return &LintResult{Valid: false}
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
	
	return lintResult
}
