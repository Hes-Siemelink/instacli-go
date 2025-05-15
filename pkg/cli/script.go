package cli

import (
	"fmt"
	"os"
)

// Script represents a CLI script to be executed
type Script struct {
	Path           string
	Debug          bool
	Output         bool
	OutputJSON     bool
	NonInteractive bool
	parsedScript   *ParsedScript
}

// NewScript creates a new Script instance
func NewScript(path string, debug, output, outputJSON, nonInteractive bool) *Script {
	return &Script{
		Path:           path,
		Debug:          debug,
		Output:         output,
		OutputJSON:     outputJSON,
		NonInteractive: nonInteractive,
	}
}

// Execute runs the script
func (s *Script) Execute() error {
	// Check if path exists
	info, err := os.Stat(s.Path)
	if err != nil {
		return fmt.Errorf("error accessing path: %w", err)
	}

	if info.IsDir() {
		return s.handleDirectory()
	}
	return s.handleFile()
}

func (s *Script) handleDirectory() error {
	// TODO: Implement directory handling
	// This should list available commands in the directory
	return fmt.Errorf("directory handling not implemented yet")
}

func (s *Script) handleFile() error {
	// Read the script file
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return fmt.Errorf("error reading script file: %w", err)
	}

	// Parse the script
	script, err := ParseScript(data)
	if err != nil {
		return fmt.Errorf("error parsing script: %w", err)
	}
	s.parsedScript = script

	// Execute the script with empty input for now
	return ExecuteScript(script, nil)
}

// GetScriptHelp returns help information for a script
func (s *Script) GetScriptHelp() (string, error) {
	if s.parsedScript == nil {
		// Read and parse the script if not already done
		data, err := os.ReadFile(s.Path)
		if err != nil {
			return "", fmt.Errorf("error reading script file: %w", err)
		}

		script, err := ParseScript(data)
		if err != nil {
			return "", fmt.Errorf("error parsing script: %w", err)
		}
		s.parsedScript = script
	}

	return GetScriptHelp(s.parsedScript), nil
}
