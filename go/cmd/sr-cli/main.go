package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Command  string
	Verbose  bool
	ShowHelp bool
}

type CommandArg struct {
	Command     string
	Description string
}

func main() {

	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	cfg := parseArgs(os.Args)

	if cfg.ShowHelp == true {
		showUsage()
		os.Exit(2)
	}

	if cfg.Command == "" {
		fmt.Println("A command is required.")
		showUsage()
		os.Exit(2)
	}

	fmt.Println(cfg.Command)
	if cfg.Verbose {
		fmt.Println("Verbose")
	} else {
		fmt.Println("Silent")
	}
}

func showUsage() {
	cmds := []CommandArg{
		{"init", "The command to execute"},
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

	var executableName = filepath.Base(os.Args[0])
	fmt.Printf("%s\t\tSwing trading toolbox\n", executableName)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Printf("%s %s\n", executableName, strings.Join(allArgs, " "))
	fmt.Println()

	for i := 0; i < len(cmds); i++ {
		fmt.Printf("%s\n", fmt.Sprintf("%-*s\t%s", colSize, cmds[i].Command, cmds[i].Description))
	}
}

func parseArgs(args []string) Config {
	cfg := Config{
		Command: "",
		Verbose: false,
	}

	for i := 1; i < len(args); i++ {
		arg := strings.ToLower(os.Args[i])
		switch arg {
		case "init":
			cfg.Command = arg
		case "-v", "--verbose":
			cfg.Verbose = true
		case "-h", "help", "--help", "-?", "?":
			cfg.ShowHelp = true
		default:
			fmt.Printf("Unknown argument: %s\n", os.Args[i])
			os.Exit(-1)
		}
	}

	return cfg
}
