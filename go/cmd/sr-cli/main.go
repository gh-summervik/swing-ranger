package main

import (
	"context"
	"database/sql"
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

	secrets, err := config.LoadSecrets()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	cfg.Secrets = secrets

	if cfg.Command == "collect" {
		cmdDb, err := sql.Open("postgres", cfg.Secrets.ConnectionStrings["Command"])
		if err != nil {
			fmt.Println("No query string found for Command.")
			os.Exit(4)
		}
		qryDb, err := sql.Open("postgres", cfg.Secrets.ConnectionStrings["Query"])
		if err != nil {
			fmt.Println("No query string found for Query.")
			os.Exit(4)
		}
		defer cmdDb.Close()
		defer qryDb.Close()

		yahooSvc := service.NewYahooService()

		ctx := context.Background()
		start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Now().UTC()

		prices, err := yahooSvc.GetHistorical(ctx, cfg.Data, start, end)
		if err != nil {
			fmt.Printf("Failed to fetch data from Yahoo: %v\n", err)
			os.Exit(6)
		}

		err = service.UpsertEodPrices(cmdDb, prices, "system")
		if err != nil {
			fmt.Printf("Failed to save data to database: %v\n", err)
			os.Exit(7)
		}

		fmt.Printf("Collected and saved %d records for %s\n", len(prices), cfg.Data)

		// chart, err := service.GetEodPrices(qryDb, cfg.Data)
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(5)
		// }
		// fmt.Println(chart)
		// return
	}

	// fmt.Printf("%+v\n", cfg.Secrets.ConnectionStrings)
	os.Exit(0)
}

func showUsage() {
	cmds := []CommandArg{
		{"[init]", "Initialize for first-time use. Will NOT injure existig initializations."},
		{"[collect <SYMBOL>]", "Capture price history for provided symbol."},
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

func parseArgs(args []string) config.Config {
	cfg := config.Config{
		Command: "",
		Verbose: false,
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
				fmt.Printf("Expecting symbol after %s\n", os.Args[i-1])
				os.Exit(-1)
			}
			cfg.Data = strings.ToUpper(os.Args[i])
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
