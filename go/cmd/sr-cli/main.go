package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/summervik/swing-ranger/internal/config"
	"github.com/summervik/swing-ranger/internal/service"

	_ "github.com/lib/pq"
)

type CommandArg struct {
	Command     string
	Description string
}

func main() {

	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	cfg, err := parseArgs(os.Args)

	if err != nil {
		fmt.Println(err)
		fmt.Println()
		showUsage()
		os.Exit(1)
	}

	// If the user asked for help, just provide the help and exit.
	if cfg.ShowHelp {
		showUsage()
		os.Exit(0)
	}

	// The absense of a command ends the process.
	if strings.TrimSpace(cfg.Command) == "" {
		fmt.Println("A command is required.")
		showUsage()
		os.Exit(2)
	}

	// Extract the secrets from the super secret json file.
	secrets, err := config.LoadSecrets()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	cfg.Secrets = secrets

	// establish services.
	commsSvc := service.NewCommsService(cfg)
	dbSvc, err := service.NewDbService(cfg, commsSvc)
	if err != nil {
		commsSvc.Communicate(fmt.Sprintf("Failed to create database service: %v", err))
		os.Exit(4)
	}

	defer dbSvc.Command.Close()
	defer dbSvc.Query.Close()

	// process the specified command.
	if cfg.Command == "collect" {
		collectSymbols(cfg, commsSvc, dbSvc)
	}

	os.Exit(0)
}

func collectSymbols(cfg config.Config, comms *service.CommsService, db *service.DbService) error {
	comms.Communicate(fmt.Sprintf("Collecting historical data for %d symbol(s)", len(cfg.Data)))
	yahooSvc := service.NewYahooService()

	ctx := context.Background()
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now().UTC()

	for _, symbol := range cfg.Data {
		comms.Communicate(fmt.Sprintf("Fetching historical data for %s", symbol))
		prices, err := yahooSvc.GetHistorical(ctx, symbol, start, end)
		if err != nil {
			comms.Communicate(fmt.Sprintf("Failed to fetch data from Yahoo for %s: %v", symbol, err))
			continue
		}

		err = db.UpsertEodPrices(prices, "system")
		if err != nil {
			return err
		}

		comms.Communicate(fmt.Sprintf("Collected and saved %d records for %s", len(prices), symbol))
	}
	return nil
}

func showUsage() {
	cmds := []CommandArg{
		{"[init]", "Initialize for first-time use. Will NOT injure existing initializations."},
		{"[collect <S1,S2,S3>]", "Capture price history for provided symbols. Use commas to separate symbols; no spaces."},
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
	fmt.Printf("%s\t\tSwing trading toolbox\n", exeName)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Printf("%s %s\n", exeName, strings.Join(allArgs, " "))
	fmt.Println()

	for i := 0; i < len(cmds); i++ {
		fmt.Printf("%s\n", fmt.Sprintf("%-*s\t%s", colSize, cmds[i].Command, cmds[i].Description))
	}
	fmt.Println()
	fmt.Println("Only one command (e.g., init, collect, etc.) is allowed; last one wins.")
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
}

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
		case "init":
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
