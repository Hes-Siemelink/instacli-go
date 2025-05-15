package cli

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Option represents a command-line option configuration
type Option struct {
	Description string `yaml:"description"`
	Default     bool   `yaml:"default,omitempty"`
	Type        string `yaml:"type,omitempty"`
	ShortOption string `yaml:"short option,omitempty"`
}

// Options represents the command-line options configuration
type Options map[string]Option

// LoadOptions loads the command-line options from the embedded YAML file
func LoadOptions() (Options, error) {
	data, err := GetSpecFile("cli/instacli-command-line-options.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading options file: %w", err)
	}

	var options Options
	if err := yaml.Unmarshal(data, &options); err != nil {
		return nil, fmt.Errorf("error parsing options file: %w", err)
	}

	return options, nil
}

// FormatHelp formats the help text for the options
func (o Options) FormatHelp() string {
	var help strings.Builder
	help.WriteString("Instacli -- Instantly create CLI applications with light scripting!\n\n")
	help.WriteString("Usage:\n")
	help.WriteString("   cli [global options] file | directory [command options]\n\n")
	help.WriteString("Global options:\n")

	// Sort options by name for consistent output
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		opt := o[name]
		shortOpt := ""
		if opt.ShortOption != "" {
			shortOpt = fmt.Sprintf(", -%s", opt.ShortOption)
		}
		help.WriteString(fmt.Sprintf("  --%s%s   %s\n", name, shortOpt, opt.Description))
	}

	return help.String()
}
