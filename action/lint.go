package action

import (
	"log/slog"
	"os"

	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/rulesets"
	"github.com/urfave/cli/v2"

	"github.com/liankui/openapi-cli/internal"
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
		slog.Error("failed to read file", "path", os.Args[1:], "error", err)
		return nil
	}

	v2 := internal.Openapi2{}
	lintResult, err := v2.Lint(c.Context, spec)
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
