package action

import (
	"log"
	"os"
	"strings"

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
			log.Printf("failed to read file, path=%v, error=%v", os.Args[1], err)
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

					log.Printf("violation, path=%v, method=%v, start=%v, end=%v, rule=%v\n", path, method, _result.StartNode.Line, _result.EndNode.Line, _result.Rule)
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
