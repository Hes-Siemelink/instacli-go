package variables

import (
	"fmt"
	"instacli/pkg/cli/commands"
	"regexp"
)

var asVarRegex = regexp.MustCompile(`^\$\{([^}]+)}$`)

type AsCommand struct {
	VarName string
}

func NewAsCommand(varName string) (*AsCommand, error) {
	m := asVarRegex.FindStringSubmatch(varName)
	if m == nil {
		return nil, fmt.Errorf("As: variable name must be in ${var} format")
	}
	return &AsCommand{VarName: m[1]}, nil
}

func (c *AsCommand) Execute(ctx *commands.ExecutionContext) error {
	output := ctx.GetOutput()
	if output == nil {
		return fmt.Errorf("As: output variable is empty")
	}
	ctx.SetVar(c.VarName, output)
	return nil
}
