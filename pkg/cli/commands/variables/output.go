package variables

import (
	"instacli/pkg/cli/commands"
)

// OutputCommand sets the output variable in the context
// It accepts any value and sets it as the output

type OutputCommand struct {
	Value interface{}
}

func NewOutputCommand(value interface{}) *OutputCommand {
	return &OutputCommand{Value: value}
}

func (c *OutputCommand) Execute(ctx *commands.ExecutionContext) error {
	ctx.SetOutput(c.Value)
	return nil
}
