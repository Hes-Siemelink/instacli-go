package testing

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// AssertEqualsCommand represents the "Assert equals" command
type AssertEqualsCommand struct {
	Actual   interface{} `yaml:"actual"`
	Expected interface{} `yaml:"expected"`
}

// Execute runs the Assert equals command
func (c *AssertEqualsCommand) Execute() error {
	if !reflect.DeepEqual(c.Actual, c.Expected) {
		actualYAML, _ := yaml.Marshal(c.Actual)
		expectedYAML, _ := yaml.Marshal(c.Expected)
		// Trim any trailing newlines from the YAML output
		actualStr := strings.TrimSpace(string(actualYAML))
		expectedStr := strings.TrimSpace(string(expectedYAML))
		return fmt.Errorf("Not equal:\n  Expected: %s\n  Actual:   %s", expectedStr, actualStr)
	}
	return nil
}

// NewAssertEquals creates a new Assert equals command
func NewAssertEquals(data map[string]interface{}) (*AssertEqualsCommand, error) {
	cmd := &AssertEqualsCommand{}

	// Extract actual and expected values
	if actual, ok := data["actual"]; ok {
		cmd.Actual = actual
	} else {
		return nil, fmt.Errorf("missing required parameter: actual")
	}

	if expected, ok := data["expected"]; ok {
		cmd.Expected = expected
	} else {
		return nil, fmt.Errorf("missing required parameter: expected")
	}

	return cmd, nil
}
