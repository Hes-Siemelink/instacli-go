package variables

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var variableRegex = regexp.MustCompile(`\$\{([^}]+)}`)

// ResolveVariablesInText replaces ${var} and ${var.path} in a string using the provided variable map.
func ResolveVariablesInText(raw string, vars map[string]interface{}) (string, error) {
	hadErr := false
	var firstErr error
	replaced := variableRegex.ReplaceAllStringFunc(raw, func(match string) string {
		varName := variableRegex.FindStringSubmatch(match)[1]
		val, err := GetValue(varName, vars)
		if err != nil {
			hadErr = true
			if firstErr == nil {
				firstErr = err
			}
			return match // leave the variable as-is
		}
		// If value is a slice or map, marshal to YAML block style
		switch v := val.(type) {
		case []interface{}, map[string]interface{}:
			yamlBytes, err := yaml.Marshal(v)
			if err != nil {
				return fmt.Sprintf("%v", v)
			}
			// Remove trailing newline for inline use
			return strings.TrimRight(string(yamlBytes), "\n")
		default:
			return fmt.Sprintf("%v", v)
		}
	})
	if hadErr {
		return replaced, firstErr
	}
	return replaced, nil
}

// GetValue resolves a variable name (with optional path) from the variable map.
func GetValue(varName string, vars map[string]interface{}) (interface{}, error) {
	name, path := splitIntoVariableAndPath(varName)
	val, ok := vars[name]
	if !ok {
		return nil, fmt.Errorf("Unknown variable ${%s}", name)
	}
	if path == "" {
		return val, nil
	}
	return resolvePath(val, path)
}

// splitIntoVariableAndPath splits 'var.path[0].foo' into ('var', '.path[0].foo')
func splitIntoVariableAndPath(varName string) (string, string) {
	for i, c := range varName {
		if c == '.' || c == '[' {
			return varName[:i], varName[i:]
		}
	}
	return varName, ""
}

// resolvePath navigates through maps/slices according to a path like '.foo[0].bar'
func resolvePath(val interface{}, path string) (interface{}, error) {
	p := path
	for len(p) > 0 {
		if p[0] == '.' {
			// Object key
			p = p[1:]
			end := strings.IndexAny(p, ".[[]")
			if end == -1 {
				end = len(p)
			}
			key := p[:end]
			m, ok := val.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Cannot access key '%s' on non-object", key)
			}
			val, ok = m[key]
			if !ok {
				return nil, fmt.Errorf("Key '%s' not found", key)
			}
			p = p[end:]
		} else if p[0] == '[' {
			// Array index
			end := strings.Index(p, "]")
			if end == -1 {
				return nil, fmt.Errorf("Unmatched '[' in path")
			}
			idxStr := p[1:end]
			var idx int
			fmt.Sscanf(idxStr, "%d", &idx)
			slice, ok := val.([]interface{})
			if !ok {
				return nil, fmt.Errorf("Cannot index non-array")
			}
			if idx < 0 || idx >= len(slice) {
				return nil, fmt.Errorf("Index %d out of range", idx)
			}
			val = slice[idx]
			p = p[end+1:]
		} else {
			return nil, fmt.Errorf("Invalid path syntax: %s", p)
		}
	}
	return val, nil
}

// ResolveVariablesRecursive recursively resolves variables in any value (string, map, slice, etc).
func ResolveVariablesRecursive(val interface{}, vars map[string]interface{}) (interface{}, error) {
	// Handle nil
	if val == nil {
		return nil, nil
	}

	switch v := val.(type) {
	case string:
		// If the string is exactly a variable reference, return the value as-is (recursively resolve if needed)
		if matches := variableRegex.FindStringSubmatch(v); matches != nil && v == matches[0] {
			resolved, err := GetValue(matches[1], vars)
			if err != nil {
				// If not found, return the original string
				return v, nil
			}
			// If the resolved value is a string that is itself a variable reference, resolve recursively
			if strVal, ok := resolved.(string); ok {
				return ResolveVariablesRecursive(strVal, vars)
			}
			return resolved, nil
		}
		return ResolveVariablesInText(v, vars)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, elem := range v {
			resolved, err := ResolveVariablesRecursive(elem, vars)
			if err != nil {
				return nil, err
			}
			result[i] = resolved
		}
		return result, nil
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, elem := range v {
			resolved, err := ResolveVariablesRecursive(elem, vars)
			if err != nil {
				return nil, err
			}
			result[key] = resolved
		}
		return result, nil
	default:
		// For other types (numbers, bools, etc), return as-is
		return val, nil
	}
}
