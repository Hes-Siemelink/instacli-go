package testing

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// AssertThatCommand represents the "Assert that" command
type AssertThatCommand struct {
	Item   interface{}   `yaml:"item,omitempty"`
	Equals interface{}   `yaml:"equals,omitempty"`
	In     interface{}   `yaml:"in,omitempty"`
	Empty  interface{}   `yaml:"empty,omitempty"`
	All    []interface{} `yaml:"all,omitempty"`
	Any    []interface{} `yaml:"any,omitempty"`
	Not    interface{}   `yaml:"not,omitempty"`
}

// Execute runs the Assert that command
func (c *AssertThatCommand) Execute() error {
	// Check for empty condition
	if c.Empty != nil {
		if !isEmpty(c.Empty) {
			emptyYAML, _ := yaml.Marshal(c.Empty)
			emptyStr := strings.TrimSpace(string(emptyYAML))
			if strings.HasPrefix(emptyStr, "-") {
				// Convert YAML list format to array format
				emptyStr = "[" + strings.Join(strings.Split(strings.TrimSpace(strings.ReplaceAll(emptyStr, "- ", "")), "\n"), ", ") + "]"
			}
			return fmt.Errorf("Condition is false.\nEmpty: %s", emptyStr)
		}
		return nil
	}

	// Check for equals condition
	if c.Item != nil && c.Equals != nil {
		if !reflect.DeepEqual(c.Item, c.Equals) {
			itemYAML, _ := yaml.Marshal(c.Item)
			equalsYAML, _ := yaml.Marshal(c.Equals)
			return fmt.Errorf("Condition is false.\nItem: %s\nEquals: %s",
				strings.TrimSpace(string(itemYAML)),
				strings.TrimSpace(string(equalsYAML)))
		}
		return nil
	}

	// Check for contains condition
	if c.Item != nil && c.In != nil {
		if !contains(c.In, c.Item) {
			itemYAML, _ := yaml.Marshal(c.Item)
			inYAML, _ := yaml.Marshal(c.In)
			// Format the error message to match the expected output
			inStr := strings.TrimSpace(string(inYAML))
			if strings.HasPrefix(inStr, "-") {
				// Convert YAML list format to array format
				inStr = "[" + strings.Join(strings.Split(strings.TrimSpace(strings.ReplaceAll(inStr, "- ", "")), "\n"), ", ") + "]"
			}
			return fmt.Errorf("Condition is false.\nItem: %s\nIn: %s",
				strings.TrimSpace(string(itemYAML)),
				inStr)
		}
		return nil
	}

	// Check for all conditions
	if c.All != nil {
		for _, condition := range c.All {
			cmd, err := NewAssertThat(condition.(map[string]interface{}))
			if err != nil {
				return err
			}
			if err := cmd.Execute(); err != nil {
				return err
			}
		}
		return nil
	}

	// Check for any conditions
	if c.Any != nil {
		var lastErr error
		for _, condition := range c.Any {
			cmd, err := NewAssertThat(condition.(map[string]interface{}))
			if err != nil {
				lastErr = err
				continue
			}
			if err := cmd.Execute(); err != nil {
				lastErr = err
				continue
			}
			return nil // At least one condition passed
		}
		return fmt.Errorf("All conditions failed. Last error: %v", lastErr)
	}

	// Check for not condition
	if c.Not != nil {
		cmd, err := NewAssertThat(c.Not.(map[string]interface{}))
		if err != nil {
			return err
		}
		if err := cmd.Execute(); err == nil {
			return fmt.Errorf("Condition is true when it should be false.\nNot: %v", c.Not)
		}
		return nil
	}

	return fmt.Errorf("no valid condition specified")
}

// isEmpty checks if a value is empty
func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String:
		return val.Len() == 0
	case reflect.Slice, reflect.Array, reflect.Map:
		return val.Len() == 0
	default:
		return false
	}
}

// contains checks if a value is contained in another value
func contains(container, item interface{}) bool {
	// If container is a map, check if item is a subset
	if reflect.TypeOf(container).Kind() == reflect.Map {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return false
		}
		containerMap := container.(map[string]interface{})
		for k, v := range itemMap {
			containerVal, exists := containerMap[k]
			if !exists || !reflect.DeepEqual(containerVal, v) {
				return false
			}
		}
		return true
	}

	// If container is a slice/array, check if item is an element
	if reflect.TypeOf(container).Kind() == reflect.Slice || reflect.TypeOf(container).Kind() == reflect.Array {
		containerVal := reflect.ValueOf(container)
		for i := 0; i < containerVal.Len(); i++ {
			if reflect.DeepEqual(containerVal.Index(i).Interface(), item) {
				return true
			}
		}
		return false
	}

	return false
}

// NewAssertThat creates a new Assert that command
func NewAssertThat(data map[string]interface{}) (*AssertThatCommand, error) {
	cmd := &AssertThatCommand{}

	// Extract conditions
	if item, ok := data["item"]; ok {
		cmd.Item = item
	}
	if equals, ok := data["equals"]; ok {
		cmd.Equals = equals
	}
	if in, ok := data["in"]; ok {
		cmd.In = in
	}
	if empty, ok := data["empty"]; ok {
		cmd.Empty = empty
	}
	if all, ok := data["all"]; ok {
		if allSlice, ok := all.([]interface{}); ok {
			cmd.All = allSlice
		} else {
			return nil, fmt.Errorf("'all' must be a list of conditions")
		}
	}
	if any, ok := data["any"]; ok {
		if anySlice, ok := any.([]interface{}); ok {
			cmd.Any = anySlice
		} else {
			return nil, fmt.Errorf("'any' must be a list of conditions")
		}
	}
	if not, ok := data["not"]; ok {
		cmd.Not = not
	}

	return cmd, nil
}
