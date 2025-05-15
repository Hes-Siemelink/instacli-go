package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runSpecFile runs all test cases in a given Instacli spec file.
func runSpecFile(t *testing.T, relPath string) {
	specRoot := os.Getenv("INSTACLI_SPEC")
	if specRoot == "" {
		t.Fatal("INSTACLI_SPEC environment variable not set")
	}
	filename := filepath.Join(specRoot, relPath)
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Split on '---' to get each test case as a full YAML document
	sections := strings.Split(string(data), "---")
	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}
		// Use the first non-empty line as the test name if it starts with 'Test case:'
		lines := strings.SplitN(section, "\n", 2)
		name := "(unnamed)"
		if len(lines) > 0 && strings.HasPrefix(lines[0], "Test case:") {
			name = strings.TrimSpace(strings.TrimPrefix(lines[0], "Test case:"))
		}
		t.Run(name, func(t *testing.T) {
			// Parse the entire section as a YAML script
			script, err := ParseScript([]byte(section))
			if err != nil {
				t.Fatalf("ParseScript error: %v", err)
			}
			if err := ExecuteScript(script, nil); err != nil {
				t.Errorf("Script execution error: %v", err)
			}
		})
	}
}

func TestInstacliSpecFiles(t *testing.T) {
	specFiles := []string{
		// Add more spec files here as needed, relative to $INSTACLI_SPEC
		"commands/instacli/variables/tests/Output variable tests.cli",
		"commands/instacli/variables/tests/Assignment tests.cli",
		"commands/instacli/variables/tests/Variable replacement tests.cli",
		// e.g. "commands/instacli/variables/tests/Other variable tests.cli",
		// e.g. "commands/instacli/db/tests/Some db tests.cli",
	}
	for _, relPath := range specFiles {
		t.Run(relPath, func(t *testing.T) {
			runSpecFile(t, relPath)
		})
	}
}
