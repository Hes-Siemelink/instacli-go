package main

import (
	"flag"
	"fmt"
	"os"

	"instacli/pkg/cli"
)

var (
	help           bool
	output         bool
	outputJSON     bool
	nonInteractive bool
	debug          bool
	options        cli.Options
)

func init() {
	// Load options from embedded YAML file
	var err error
	options, err = cli.LoadOptions()
	if err != nil {
		fmt.Printf("Error loading options: %v\n", err)
		os.Exit(1)
	}

	// Register flags based on options
	for name, opt := range options {
		switch name {
		case "help":
			flag.BoolVar(&help, name, opt.Default, opt.Description)
			if opt.ShortOption != "" {
				flag.BoolVar(&help, opt.ShortOption, opt.Default, opt.Description)
			}
		case "output":
			flag.BoolVar(&output, name, opt.Default, opt.Description)
			if opt.ShortOption != "" {
				flag.BoolVar(&output, opt.ShortOption, opt.Default, opt.Description)
			}
		case "output-json":
			flag.BoolVar(&outputJSON, name, opt.Default, opt.Description)
			if opt.ShortOption != "" {
				flag.BoolVar(&outputJSON, opt.ShortOption, opt.Default, opt.Description)
			}
		case "non-interactive":
			flag.BoolVar(&nonInteractive, name, opt.Default, opt.Description)
			if opt.ShortOption != "" {
				flag.BoolVar(&nonInteractive, opt.ShortOption, opt.Default, opt.Description)
			}
		case "debug":
			flag.BoolVar(&debug, name, opt.Default, opt.Description)
			if opt.ShortOption != "" {
				flag.BoolVar(&debug, opt.ShortOption, opt.Default, opt.Description)
			}
		}
	}
}

func printUsage() {
	fmt.Print(options.FormatHelp())
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 || help {
		printUsage()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error: No script or directory specified")
		os.Exit(1)
	}

	script := cli.NewScript(args[0], debug, output, outputJSON, nonInteractive)

	if help {
		helpText, err := script.GetScriptHelp()
		if err != nil {
			fmt.Printf("Error getting help: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(helpText)
		return
	}

	if err := script.Execute(); err != nil {
		if debug {
			fmt.Printf("Error: %+v\n", err)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		os.Exit(1)
	}
}
