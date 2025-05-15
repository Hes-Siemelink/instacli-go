package testing

import (
	"fmt"
	"instacli/pkg/cli/commands"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type ExpectedOutputCommand struct {
	Expected interface{}
}

func NewExpectedOutputCommand(expected interface{}) *ExpectedOutputCommand {
	return &ExpectedOutputCommand{Expected: expected}
}

func (c *ExpectedOutputCommand) Execute(ctx *commands.ExecutionContext) error {
	actual := ctx.GetOutput()
	if !reflect.DeepEqual(actual, c.Expected) {
		actualYAML, _ := yaml.Marshal(actual)
		expectedYAML, _ := yaml.Marshal(c.Expected)
		return fmt.Errorf("Unexpected output.\n  Expected: %s\n  Actual:   %s", strings.TrimSpace(string(expectedYAML)), strings.TrimSpace(string(actualYAML)))
	}
	return nil
}
