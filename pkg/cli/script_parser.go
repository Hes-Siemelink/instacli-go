package cli

import (
	"fmt"
	"reflect"
	"strings"

	"instacli/pkg/cli/commands"
	"instacli/pkg/cli/commands/testing"
	"instacli/pkg/cli/commands/variables"

	"bytes"

	"gopkg.in/yaml.v3"
)

// ScriptMetadata represents the metadata section of a script
type ScriptMetadata struct {
	Description string                `yaml:"description"`
	Input       map[string]InputParam `yaml:"input"`
}

// InputParam represents an input parameter definition
type InputParam struct {
	Description string `yaml:"description"`
	Default     string `yaml:"default,omitempty"`
}

// ScriptCommand represents a single command in the script
type ScriptCommand struct {
	Print          string                 `yaml:"Print,omitempty"`
	AssertEquals   map[string]interface{} `yaml:"Assert equals,omitempty"`
	AssertThat     map[string]interface{} `yaml:"Assert that,omitempty"`
	Output         interface{}            `yaml:"Output,omitempty"`
	ExpectedOutput interface{}            `yaml:"Expected output,omitempty"`
	As             string                 `yaml:"As,omitempty"`
	TestCase       string                 `yaml:"Test case,omitempty"`
	VarAssign      *VarAssignment         `yaml:"-"`
}

type VarAssignment struct {
	Name  string
	Value interface{}
}

// ParsedScript represents a parsed Instacli script
type ParsedScript struct {
	Metadata ScriptMetadata
	Commands []ScriptCommand
}

// ParseScript parses a script file into a ParsedScript struct
func ParseScript(data []byte) (*ParsedScript, error) {
	content := string(data)
	script := &ParsedScript{
		Metadata: ScriptMetadata{
			Input: make(map[string]InputParam),
		},
	}

	// varAssignRegex := regexp.MustCompile(`(?m)^\$\{([^}]+)}:\s*(.*)$`)

	// Check if the script has metadata (contains "---" separator)
	if strings.Contains(content, "---\n") {
		parts := strings.SplitN(content, "---\n", 2)

		// Parse metadata section
		if err := yaml.Unmarshal([]byte(parts[0]), &script.Metadata); err != nil {
			return nil, fmt.Errorf("error parsing script info: %w", err)
		}

		// Parse commands section
		if err := yaml.Unmarshal([]byte(parts[1]), &script.Commands); err != nil {
			return nil, fmt.Errorf("error parsing commands: %w", err)
		}
	} else {
		// Use yaml.Decoder to parse multiple YAML documents or multiple top-level keys
		dec := yaml.NewDecoder(bytes.NewReader(data))
		for {
			var m map[string]interface{}
			if err := dec.Decode(&m); err != nil {
				if err.Error() == "EOF" {
					break
				}
				return nil, fmt.Errorf("error decoding YAML: %w", err)
			}
			if len(m) == 0 {
				continue
			}
			// For each top-level key, create a ScriptCommand
			for k, v := range m {
				cmd := ScriptCommand{}
				switch k {
				case "Print":
					cmd.Print, _ = v.(string)
				case "Assert equals":
					if mm, ok := v.(map[string]interface{}); ok {
						cmd.AssertEquals = mm
					}
				case "Assert that":
					if mm, ok := v.(map[string]interface{}); ok {
						cmd.AssertThat = mm
					}
				case "Output":
					cmd.Output = v
				case "Expected output":
					cmd.ExpectedOutput = v
				case "As":
					cmd.As, _ = v.(string)
				case "Test case":
					cmd.TestCase, _ = v.(string)
				default:
					// Try variable assignment
					if strings.HasPrefix(k, "${") && strings.HasSuffix(k, "}") {
						cmd.VarAssign = &VarAssignment{Name: k[2 : len(k)-1], Value: v}
					}
				}
				script.Commands = append(script.Commands, cmd)
			}
		}
	}

	return script, nil
}

// ExecuteScript runs the script with the given input parameters
func ExecuteScript(script *ParsedScript, input map[string]string) error {
	ctx := commands.NewExecutionContext()

	// First pass: fixed-point iteration for variable assignments
	maxIterations := 10
	for i := 0; i < maxIterations; i++ {
		progress := false
		for _, cmd := range script.Commands {
			if cmd.VarAssign != nil {
				val := cmd.VarAssign.Value
				vars := ctx.Vars()
				_, alreadySet := vars[cmd.VarAssign.Name]
				resolved, err := variables.ResolveVariablesRecursive(val, vars)
				if err == nil {
					if !alreadySet || !reflect.DeepEqual(vars[cmd.VarAssign.Name], resolved) {
						ctx.SetVar(cmd.VarAssign.Name, resolved)
						progress = true
					}
				}
			}
		}
		if !progress {
			break
		}
	}

	// Second pass: process all other commands (skip only true variable assignments)
	for _, cmd := range script.Commands {
		if cmd.TestCase != "" {
			continue
		}
		if cmd.VarAssign != nil {
			// Already processed in first pass
			continue
		}
		if cmd.As != "" {
			asCmd, err := variables.NewAsCommand(cmd.As)
			if err != nil {
				return fmt.Errorf("error creating As command: %w", err)
			}
			if err := asCmd.Execute(ctx); err != nil {
				return fmt.Errorf("error executing As command: %w", err)
			}
			continue
		}
		vars := ctx.Vars()
		if cmd.Print != "" {
			output, _ := variables.ResolveVariablesInText(cmd.Print, vars)
			fmt.Println(output)
			ctx.SetOutput(output)
		}
		if cmd.Output != nil {
			val, err := variables.ResolveVariablesRecursive(cmd.Output, vars)
			if err != nil {
				return fmt.Errorf("error resolving variables in Output: %w", err)
			}
			outputCmd := variables.NewOutputCommand(val)
			if err := outputCmd.Execute(ctx); err != nil {
				return fmt.Errorf("error executing Output command: %w", err)
			}
		}
		if cmd.ExpectedOutput != nil {
			expected, err := variables.ResolveVariablesRecursive(cmd.ExpectedOutput, vars)
			if err != nil {
				return fmt.Errorf("error resolving variables in Expected output: %w", err)
			}
			actual := ctx.GetOutput()
			actualYAML, _ := yaml.Marshal(actual)
			expectedYAML, _ := yaml.Marshal(expected)
			actualStr := strings.TrimSpace(string(actualYAML))
			expectedStr := strings.TrimSpace(string(expectedYAML))
			if actualStr != expectedStr {
				return fmt.Errorf("Unexpected output.\n  Expected: %s\n  Actual:   %s", expectedStr, actualStr)
			}
		}
		if cmd.AssertEquals != nil {
			resolved, err := variables.ResolveVariablesRecursive(cmd.AssertEquals, vars)
			if err != nil {
				return fmt.Errorf("error resolving variables in Assert equals: %w", err)
			}
			assertCmd, err := testing.NewAssertEquals(resolved.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("error creating Assert equals command: %w", err)
			}
			// Trim whitespace for actual and expected if they are strings
			if a, ok := assertCmd.Actual.(string); ok {
				assertCmd.Actual = strings.TrimSpace(a)
			}
			if e, ok := assertCmd.Expected.(string); ok {
				assertCmd.Expected = strings.TrimSpace(e)
			}
			if err := assertCmd.Execute(); err != nil {
				return fmt.Errorf("assertion failed: %w", err)
			}
		}
		if cmd.AssertThat != nil {
			resolved, err := variables.ResolveVariablesRecursive(cmd.AssertThat, vars)
			if err != nil {
				return fmt.Errorf("error resolving variables in Assert that: %w", err)
			}
			assertCmd, err := testing.NewAssertThat(resolved.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("error creating Assert that command: %w", err)
			}
			if err := assertCmd.Execute(); err != nil {
				return fmt.Errorf("assertion failed: %w", err)
			}
		}
	}
	return nil
}

// GetScriptHelp returns the help text for the script
func GetScriptHelp(script *ParsedScript) string {
	var help strings.Builder
	help.WriteString(script.Metadata.Description + "\n\n")

	if len(script.Metadata.Input) > 0 {
		help.WriteString("Options:\n")
		for name, param := range script.Metadata.Input {
			help.WriteString(fmt.Sprintf("  --%s   %s\n", name, param.Description))
		}
	}

	return help.String()
}
