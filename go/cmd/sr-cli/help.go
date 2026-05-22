package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/summervik/swing-ranger/internal/config"

	_ "github.com/lib/pq"
)

func parseArgs(args []string) (config.Config, error) {
	cfg := config.Config{
		Command:  "",
		Verbose:  false,
		ShowHelp: false,
		Data:     nil,
		Secrets:  nil,
	}

	count := len(args)

	for i := 1; i < count; i++ {
		arg := strings.ToLower(os.Args[i])
		switch arg {
		case "init", "update", "test":
			cfg.Command = arg
		case "collect":
			cfg.Command = arg
			i++
			if i >= count {
				return cfg, fmt.Errorf("Expecting one or more symbols after %s", args[i-1])
			}

			symbolsStr := strings.TrimSpace(args[i])
			if symbolsStr == "" {
				return cfg, fmt.Errorf("Expecting symbol(s) after %s\n", args[i-1])
			}
			rawSymbols := strings.Split(symbolsStr, ",")
			cfg.Data = make([]string, 0, len(rawSymbols))
			for _, s := range rawSymbols {
				trimmed := strings.TrimSpace(s)
				if trimmed != "" {
					cfg.Data = append(cfg.Data, strings.ToUpper(trimmed))
				}
			}
			if len(cfg.Data) == 0 {
				return cfg, fmt.Errorf("No valid symbols provided")
			}
		case "-v", "--verbose":
			cfg.Verbose = true
		case "-h", "help", "--help", "-?", "?":
			cfg.ShowHelp = true
		default:
			return cfg, fmt.Errorf("Unknown argument: %s\n", os.Args[i])
		}
	}

	return cfg, nil
}

func showUsage() {
	cmds := []CommandArg{
		{"[init]", "Command to initialize first-time use. Will NOT injure existing initializations."},
		{"[collect <S1,S2,S3>]", "Command to capture price history for provided symbols. Use commas to separate symbols; no spaces."},
		{"[update]", "Command to update price action data for existing symbols (symbols previously captured)."},
		{"[-h | -? | ? | --help]", "Show help"},
		{"[-v | --verbose]", "Verbose output"},
	}

	allArgs := make([]string, len(cmds))
	for i, cmd := range cmds {
		allArgs[i] = cmd.Command
	}

	var colSize = 0
	for i := 0; i < len(cmds); i++ {

		if len(cmds[i].Command) > colSize {
			colSize = len(cmds[i].Command)
		}
	}

	var exeName = filepath.Base(os.Args[0])
	fmt.Printf("\n%s\t\tSwing trading toolbox\n", exeName)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Printf("%s %s\n", exeName, strings.Join(allArgs, " "))
	fmt.Println()

	for i := 0; i < len(cmds); i++ {
		fmt.Printf("%s\n", fmt.Sprintf("%-*s\t%s", colSize, cmds[i].Command, cmds[i].Description))
	}
	fmt.Println()
	fmt.Println("Only one command (e.g., init, collect, update, etc.) is allowed; last one wins.")
	fmt.Println()
	fmt.Println("Examples")
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println()
	fmt.Printf("%s collect MSFT,TSLA", exeName)
	fmt.Println("\tWill collect and preserve historical price information for MSFT and TSLA.")
	fmt.Println()
	fmt.Printf("%s collect MSFT,TSLA -v", exeName)
	fmt.Println("\tCollects historical price information for MSFT and TSLA with verbose output.")
	fmt.Println()
	fmt.Printf("%s update", exeName)
	fmt.Println("\tUpdates price action for existing symbols.")
	fmt.Println()
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println()
	fmt.Println("A secrets.json file is required; here is an example of the content.")
	fmt.Println()
	fmt.Println(`{
  "ConnectionStrings": {
    "Command": "Your connection string here",
    "Query": "Your connection string here"
  }
}`)
}
