package action

import (
	"os"
	"strings"

	"github.com/chaos-io/chaos/core/logs"
	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/daveshanley/vacuum/rulesets"
	"github.com/urfave/cli/v2"
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

func (s *Lint) Action() func(c *cli.Context) error {
	return func(c *cli.Context) error {

		spec, err := os.ReadFile(c.Args().First())
		if err != nil {
			logs.Errorw("failed to read file", "error", err, "path", os.Args[1])
			return err
		}

		result := motor.ApplyRulesToRuleSet(&motor.RuleSetExecution{
			RuleSet:         lintRules,
			Spec:            spec,
			CustomFunctions: map[string]model.RuleFunction{},
		})

		operations := result.Index.GetAllParametersFromOperations()
		for path, operation := range operations {
			for method := range operation {
				operationPath := strings.Join([]string{"$.paths", path, method, "parameters"}, ".")
				for _, _result := range result.Results {
					if _result.Path != operationPath {
						continue
					}

					logs.Infow("violation", "path", path, "method", method, "start", _result.StartNode.Line, "end", _result.EndNode.Line, "rule", _result.Rule)
					break
				}

			}
		}

		return nil
	}
}

var (
	lintRules = &rulesets.RuleSet{
		Rules: map[string]*model.Rule{
			"operation-parameters": rulesets.GetOperationParametersRule(),
		},
	}
)
